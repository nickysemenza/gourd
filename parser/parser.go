// Package parser turns a human-readable ingredient line-item into a structed format.
package parser

// go:generate stringer -type=TokenKind

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/nickysemenza/gourd/unit"
	log "github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
)

type tokenKind int

// types of things we can parse
const (
	// strings
	measureWeight tokenKind = iota
	measureVolume
	unkString
	kindIngredientName
	modifier

	// floats
	unkFloat
	weightFloat
	volumeFloat

	none
)

type segment struct {
	kind tokenKind
	raw  string
}

type parser struct {
	current        tokenKind
	sb             strings.Builder
	nextIsModifier bool
	res            []segment
}

// Measurement represents a volume or weight measurement.
type Measurement struct {
	Unit  string
	Value float64
}

// Ingredient is the structured version of an ingredient line-item.
type Ingredient struct {
	Name     string
	Weight   Measurement
	Volume   Measurement
	Modifier string
}

func (i *Ingredient) Grams() float64 {
	if !unit.IsGram(i.Weight.Unit) {
		return 0
	}
	return i.Weight.Value
}

// ToString converts the structure back to a string. This will be a normalized version of what was parsed, if applicable.
func (i *Ingredient) ToString() string {
	var weight, volume string
	if i.Weight.Value != 0 {
		weight = fmt.Sprintf("%g %s", i.Weight.Value, i.Weight.Unit)
	}
	if i.Volume.Value != 0 {
		volume = fmt.Sprintf("%g %s", i.Volume.Value, i.Volume.Unit)
	}

	var sb strings.Builder
	switch {
	case weight != "" && volume != "":
		sb.WriteString(fmt.Sprintf("%s (%s)", weight, volume))
	case weight != "":
		sb.WriteString(weight)
	case volume != "":
		sb.WriteString(volume)
	default:
		sb.WriteString("0 g")
	}
	sb.WriteString(" ")
	sb.WriteString(i.Name)
	if i.Modifier != "" {
		sb.WriteString(fmt.Sprintf(", %s", i.Modifier))
	}

	return sb.String()
}

// Parse attempts to parse the line-item.
func Parse(ctx context.Context, s string) (*Ingredient, error) {
	return (&parser{}).parse(ctx, s)
}

func (p *parser) parse(ctx context.Context, s string) (*Ingredient, error) {
	ctx, span := otel.Tracer("parser").Start(ctx, "parser.Parse")
	defer span.End()

	segments, err := p.getsegments(ctx, s)
	span.AddEvent("got segments")
	span.SetAttributes(label.String("raw", s))
	if err != nil {
		return nil, fmt.Errorf("failed to get segments: %w", err)
	}
	if len(segments) < 2 { //nolint:gomnd
		return nil, fmt.Errorf("no enough segments found")
	}
	// spew.Dump(segments)
	log.Debugf("INPUT: %s\nOUTPUT:\n", s)
	return p.handleSegments(ctx, segments)
}

// nolint:gocognit,funlen,exhaustive
func (p *parser) handleSegments(ctx context.Context, segments []segment) (*Ingredient, error) {
	_, span := otel.Tracer("parser").Start(ctx, "parser.handleSegments")
	defer span.End()
	ing := Ingredient{}
	for i := 0; i < len(segments); i++ {
		curr := segments[i]
		switch curr.kind {
		case volumeFloat, weightFloat:
			f, err := parseFloat(curr.raw)
			if err != nil {
				return nil, err
			}

			// look back one and see if there is another float that needs to be summed with this one (e.g. `1`,`1/2`)
			if i > 0 {
				prev := segments[i-1]
				if prev.kind == unkFloat {
					prevFloat, err := parseFloat(prev.raw)
					if err != nil {
						return nil, err
					}
					f += prevFloat
				}
			}

			// look ahead one to find the unit that pairs with it
			if !(i < len(segments)) {
				return nil, fmt.Errorf("failed to look ahead and find matching measurement unit for %v", curr)
			}
			next := segments[i+1]
			log.Debugf("%g %s (%s, %s)\n", f, next.raw, curr.kind, next.kind)
			i++

			m := Measurement{Unit: next.raw, Value: f}
			switch next.kind {
			case measureVolume:
				ing.Volume = m
			case measureWeight:
				ing.Weight = m
			default:
				return nil, fmt.Errorf("failed to look ahead and find matching measurement unit for %v (%v), next was %s", curr, curr.kind, next.kind)
			}

		case modifier:
			log.Debugf("%s (%s)\n", curr.raw, curr.kind)
			ing.Modifier = curr.raw
		case kindIngredientName:
			// join multiple parts of the ingredient name back together
			var ings []string
			for x := i; x < len(segments); x++ {
				curr2 := segments[x]
				if curr2.kind != kindIngredientName {
					break
				}
				ings = append(ings, curr2.raw)
			}
			name := strings.Join(ings, " ")
			log.Debugf("%s (%s)\n", name, kindIngredientName)
			i += len(ings) - 1
			ing.Name = name
		}
	}
	return &ing, nil
}

// nolint:exhaustive
func (p *parser) getsegments(ctx context.Context, s string) ([]segment, error) {
	_, span := otel.Tracer("parser").Start(ctx, "parser.getsegments")
	defer span.End()
	p.current = none
	r := strings.NewReader(s)
	p.sb.Reset()
	p.res = []segment{}

	for {
		ch, _, err := r.ReadRune() // err will only ever be EOF

		switch {
		case unicode.IsSpace(ch), errors.Is(err, io.EOF):
			p.handleDone()
			if errors.Is(err, io.EOF) {
				return p.res, nil
			}
		case unicode.IsDigit(ch):
			log.Debug("found digit")
			p.current = unkFloat
			p.sb.WriteRune(ch)

		case unicode.IsNumber(ch):
			log.Debug("found number")
			p.current = unkFloat
			p.sb.WriteString(runeNumberToString(ch))

		case unicode.IsPunct(ch):
			log.Debug("found punct")
			if p.current == unkFloat {
				p.sb.WriteRune(ch)
				// period in middle of decimal
			} else {
				p.handleDone()
				if ch == ',' {
					p.nextIsModifier = true
				}
			}
		case unicode.IsLetter(ch):
			log.Debug("found letter")
			switch p.current {
			case none, unkString:
				p.current = unkString
				p.sb.WriteRune(ch)
			case unkFloat:
				p.handleDone()
				p.current = unkString
				p.sb.WriteRune(ch)
				// 1.2g
			}

		default:
			return nil, fmt.Errorf("unknown rune %v", ch)
		}
	}
}
func (p *parser) handleDone() {
	last := p.sb.String()
	lastUnkFloatShouldBe := none

	// if contains(weightUnits, last) {
	if unit.IsWeight(last) {
		p.current = measureWeight
		lastUnkFloatShouldBe = weightFloat
	}
	if unit.IsVolume(last) {
		p.current = measureVolume
		lastUnkFloatShouldBe = volumeFloat
	}
	if p.current == unkString {
		if p.nextIsModifier {
			p.current = modifier
			// p.nextIsModifier = false
		} else {
			p.current = kindIngredientName
		}
	}

	// fmt.Printf("last was: %s (%s)\n", last, p.current)
	p.res = append(p.res, segment{p.current, last})

	if lastUnkFloatShouldBe != none {
		for i := len(p.res) - 1; i >= 0; i-- {
			if p.res[i].kind == unkFloat {
				p.res[i].kind = lastUnkFloatShouldBe
				break
			}
		}
	}

	p.sb.Reset()
	p.current = none
}

//nolint:gomnd
func runeNumberToString(r rune) string {
	switch r {
	case 188:
		return "1/4"
	case 189:
		return "1/2"
	case 190:
		return "3/4"
	default:
		return ""
	}
}

func parseFloat(s string) (float64, error) {
	if strings.Contains(s, "/") {
		parts := strings.Split(s, "/")
		p1, err := parseFloat(parts[0])
		if err != nil {
			return 0, fmt.Errorf("failed to parse fractional float: %w", err)
		}
		p2, err := parseFloat(parts[1])
		if err != nil {
			return 0, fmt.Errorf("failed to parse fractional float: %w", err)
		}

		return p1 / p2, nil
	}
	return strconv.ParseFloat(s, 10)
}

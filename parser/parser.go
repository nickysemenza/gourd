//go:generate stringer -type=TokenKind

package parser

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/nickysemenza/food/unit"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/global"
)

type TokenKind int

// types of things we can parse
const (
	// strings
	MeasureWeight TokenKind = iota
	MeasureVolume
	UnkString
	KindIngredientName
	Modifier

	// floats
	UnkFloat
	WeightFloat
	VolumeFloat

	None
)

type segment struct {
	kind TokenKind
	raw  string
}

type parser struct {
	current        TokenKind
	sb             strings.Builder
	nextIsModifier bool
	res            []segment
}

type Measurement struct {
	Unit  string
	Value float64
	// kind  TokenKind
}
type Ingredient struct {
	Name     string
	Weight   Measurement
	Volume   Measurement
	Modifier string
}

func (i *Ingredient) ToString() string {
	var sb strings.Builder

	sb.WriteString(i.Name)
	sb.WriteString("\n")
	if i.Weight.Value != 0 {
		sb.WriteString(fmt.Sprintf("•%g %s\n", i.Weight.Value, i.Weight.Unit))
	}
	if i.Volume.Value != 0 {
		sb.WriteString(fmt.Sprintf("•%g %s\n", i.Volume.Value, i.Volume.Unit))
	}
	return sb.String()
}

func Parse(ctx context.Context, s string) (*Ingredient, error) {
	return (&parser{}).parse(ctx, s)
}

func (p *parser) parse(ctx context.Context, s string) (*Ingredient, error) {
	ctx, span := global.Tracer("parser").Start(ctx, "parser.Parse")
	defer span.End()

	segments, err := p.getsegments(s)
	span.AddEvent(ctx, "got segments")
	span.SetAttributes(core.KeyValue{Key: "raw", Value: core.String(s)})
	if err != nil {
		return nil, fmt.Errorf("failed to get segments: %v", err)
	}
	if len(segments) < 2 {
		return nil, fmt.Errorf("no enough segments found")
	}
	// spew.Dump(segments)
	log.Debugf("INPUT: %s\nOUTPUT:\n", s)
	return p.handleSegments(segments)
}

// nolint:gocognit,funlen
func (p *parser) handleSegments(segments []segment) (*Ingredient, error) {
	ing := Ingredient{}
	for i := 0; i < len(segments); i++ {
		curr := segments[i]
		switch curr.kind {
		case VolumeFloat, WeightFloat:
			f, err := parseFloat(curr.raw)
			if err != nil {
				return nil, err
			}

			// look back one and see if there is another float that needs to be summed with this one (e.g. `1`,`1/2`)
			if i > 0 {
				prev := segments[i-1]
				if prev.kind == UnkFloat {
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
			case MeasureVolume:
				ing.Volume = m
			case MeasureWeight:
				ing.Weight = m
			default:
				return nil, fmt.Errorf("failed to look ahead and find matching measurement unit for %v, next was %s", curr, next.kind)
			}

		case Modifier:
			log.Debugf("%s (%s)\n", curr.raw, curr.kind)
			ing.Modifier = curr.raw
		case KindIngredientName:
			// join multiple parts of the ingredient name back together
			var ings []string
			for x := i; x < len(segments); x++ {
				curr2 := segments[x]
				if curr2.kind != KindIngredientName {
					break
				}
				ings = append(ings, curr2.raw)
			}
			name := strings.Join(ings, " ")
			log.Debugf("%s (%s)\n", name, KindIngredientName)
			i += len(ings) - 1
			ing.Name = name
		}
	}
	return &ing, nil
}
func (p *parser) getsegments(s string) ([]segment, error) {
	p.current = None
	r := strings.NewReader(s)
	p.sb.Reset()
	p.res = []segment{}

	for {
		ch, _, err := r.ReadRune() // err will only ever be EOF

		switch {
		case unicode.IsSpace(ch), err == io.EOF:
			p.handleDone()
			if err == io.EOF {
				return p.res, nil
			}
		case unicode.IsDigit(ch):
			log.Debug("found digit")
			p.current = UnkFloat
			p.sb.WriteRune(ch)

		case unicode.IsNumber(ch):
			log.Debug("found number")
			p.current = UnkFloat
			p.sb.WriteString(runeNumberToString(ch))

		case unicode.IsPunct(ch):
			log.Debug("found punct")
			if p.current == UnkFloat {
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
			case None, UnkString:
				p.current = UnkString
				p.sb.WriteRune(ch)
			case UnkFloat:
				p.handleDone()
				p.current = UnkString
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
	lastUnkFloatShouldBe := None

	// if contains(weightUnits, last) {
	if unit.IsWeight(last) {
		p.current = MeasureWeight
		lastUnkFloatShouldBe = WeightFloat
	}
	if unit.IsVolume(last) {
		p.current = MeasureVolume
		lastUnkFloatShouldBe = VolumeFloat
	}
	if p.current == UnkString {
		if p.nextIsModifier {
			p.current = Modifier
			p.nextIsModifier = false
		} else {
			p.current = KindIngredientName
		}
	}

	// fmt.Printf("last was: %s (%s)\n", last, p.current)
	p.res = append(p.res, segment{p.current, last})

	if lastUnkFloatShouldBe != None {
		for i := len(p.res) - 1; i >= 0; i-- {
			if p.res[i].kind == UnkFloat {
				p.res[i].kind = lastUnkFloatShouldBe
				break
			}
		}
	}

	p.sb.Reset()
	p.current = None
}

func runeNumberToString(r rune) string {
	switch r {
	case 189:
		return "1/2"
	case 190:
		return "3/4"
	default:
		return ""
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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

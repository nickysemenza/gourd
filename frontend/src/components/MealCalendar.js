import React, { Component } from 'react';
import moment from 'moment';
import 'moment-timezone';
import BigCalendar from 'react-big-calendar';
import { List } from 'semantic-ui-react';

class MealCalendar extends Component {
  Event = ({ event }) => {
    return (
      <span>
        <strong>{event.title}</strong>
        {event.desc && ':  ' + event.desc}
      </span>
    );
  };

  EventAgenda = ({ event }) => {
    return (
      <span>
        <em>{event.title}</em>
        <p>{event.desc}</p>
        <List link>{event.extra}</List>
      </span>
    );
  };

  customSlotPropGetter = date => {
    if (date.getDate() === 7 || date.getDate() === 15)
      return {
        className: 'special-day'
      };
    else return {};
  };

  render() {
    let events = this.props.meal_list.map(meal => {
      let d = moment(meal.time)
        .startOf('day')
        .toDate();
      return {
        title: `${meal.type}: ${meal.name}`,
        startDate: d,
        endDate: d,
        desc: meal.description,
        extra: this.props.buildRecipeWithMultiplierListForMeal(meal)
      };
    });
    BigCalendar.momentLocalizer(moment); // or globalizeLocalizer

    return (
      <BigCalendar
        events={events}
        startAccessor="startDate"
        endAccessor="endDate"
        slotPropGetter={this.customSlotPropGetter}
        components={{
          event: this.Event,
          agenda: {
            event: this.EventAgenda
          }
        }}
      />
    );
  }
}

export default MealCalendar;

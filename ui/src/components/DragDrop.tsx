import React, { useRef } from "react";
import { useDrop, DropTargetMonitor, XYCoord, useDrag } from "react-dnd";

export interface DragWrapperProps {
  id: any;
  index: number;
  moveCard: (dragIndex: number, hoverIndex: number) => void;
  enable: boolean;
}

interface DragItem {
  index: number;
  id: string;
  type: string;
}

export const DragWrapper: React.FC<DragWrapperProps> = ({
  id,
  index,
  moveCard,
  children,
  enable,
}) => {
  // https://github.com/react-dnd/react-dnd/blob/main/packages/examples-hooks/src/04-sortable/simple/Card.tsx
  const ref = useRef<HTMLDivElement>(null);
  const [, drop] = useDrop({
    accept: "card1",
    hover(item: DragItem, monitor: DropTargetMonitor) {
      if (!ref.current) {
        return;
      }
      const dragIndex = item.index;
      const hoverIndex = index;

      // Don't replace items with themselves
      if (dragIndex === hoverIndex) {
        return;
      }

      // Determine rectangle on screen
      const hoverBoundingRect = ref.current?.getBoundingClientRect();

      // Get vertical middle
      const hoverMiddleY =
        (hoverBoundingRect.bottom - hoverBoundingRect.top) / 2;

      // Determine mouse position
      const clientOffset = monitor.getClientOffset();

      // Get pixels to the top
      const hoverClientY = (clientOffset as XYCoord).y - hoverBoundingRect.top;

      // Only perform the move when the mouse has crossed half of the items height
      // When dragging downwards, only move when the cursor is below 50%
      // When dragging upwards, only move when the cursor is above 50%

      // Dragging downwards
      if (dragIndex < hoverIndex && hoverClientY < hoverMiddleY) {
        return;
      }

      // Dragging upwards
      if (dragIndex > hoverIndex && hoverClientY > hoverMiddleY) {
        return;
      }

      // Time to actually perform the action
      moveCard(dragIndex, hoverIndex);

      // Note: we're mutating the monitor item here!
      // Generally it's better to avoid mutations,
      // but it's good here for the sake of performance
      // to avoid expensive index searches.
      item.index = hoverIndex;
    },
  });

  const [{ isDragging }, drag] = useDrag({
    item: { type: "card1", id, index },
    collect: (monitor: any) => ({
      isDragging: monitor.isDragging(),
    }),
  });

  const opacity = isDragging ? 0 : 1;
  if (enable) drag(drop(ref));
  return (
    <div ref={ref} style={{ opacity }}>
      {children}
    </div>
  );
};

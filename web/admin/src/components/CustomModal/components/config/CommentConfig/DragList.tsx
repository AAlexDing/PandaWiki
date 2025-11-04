import {
  closestCenter,
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
  MouseSensor,
  TouchSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import {
  arrayMove,
  rectSortingStrategy,
  SortableContext,
} from '@dnd-kit/sortable';
import { Stack } from '@mui/material';
import { Dispatch, FC, SetStateAction, useCallback, useState } from 'react';
import Item, { type ItemType } from './Item';
import SortableItem from './SortableItem';

interface DragListProps {
  data: ItemType[];
  onChange: (data: ItemType[]) => void;
  setIsEdit: Dispatch<SetStateAction<boolean>>;
}

const DragList: FC<DragListProps> = ({ data, onChange, setIsEdit }) => {
  const [activeId, setActiveId] = useState<string | null>(null);
  const sensors = useSensors(useSensor(MouseSensor), useSensor(TouchSensor));

  const handleDragStart = useCallback((event: DragStartEvent) => {
    setActiveId(event.active.id as string);
  }, []);

  const handleDragEnd = useCallback(
    (event: DragEndEvent) => {
      const { active, over } = event;
      if (active.id !== over?.id) {
        const oldIndex = data.findIndex(item => item.id === active.id);
        const newIndex = data.findIndex(item => item.id === over!.id);
        const newData = arrayMove(data, oldIndex, newIndex);
        onChange(newData);
      }
      setActiveId(null);
    },
    [data, onChange],
  );

  const handleDragCancel = useCallback(() => {
    setActiveId(null);
  }, []);

  const handleRemove = useCallback(
    (id: string) => {
      const newData = data.filter(item => item.id !== id);
      onChange(newData);
    },
    [data, onChange],
  );

  const handleUpdateItem = useCallback(
    (updatedItem: ItemType) => {
      const newData = data.map(item =>
        item.id === updatedItem.id ? updatedItem : item,
      );
      onChange(newData);
    },
    [data, onChange],
  );

  if (data.length === 0) return null;

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCenter}
      onDragStart={handleDragStart}
      onDragEnd={handleDragEnd}
      onDragCancel={handleDragCancel}
    >
      <SortableContext
        items={data.map(item => item.id)}
        strategy={rectSortingStrategy}
      >
        <Stack direction={'row'} flexWrap={'wrap'} gap={2}>
          {data.map(item => (
            <SortableItem
              key={item.id}
              id={item.id}
              item={item}
              handleRemove={handleRemove}
              handleUpdateItem={handleUpdateItem}
              setIsEdit={setIsEdit}
            />
          ))}
        </Stack>
      </SortableContext>
      <DragOverlay adjustScale style={{ transformOrigin: '0 0' }}>
        {activeId ? (
          <Item
            isDragging
            item={data.find(item => item.id === activeId)!}
            setIsEdit={setIsEdit}
            handleUpdateItem={handleUpdateItem}
          />
        ) : null}
      </DragOverlay>
    </DndContext>
  );
};

export default DragList;

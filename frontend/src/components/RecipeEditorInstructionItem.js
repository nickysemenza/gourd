import React from 'react';
import { Button, Grid, Input, Label } from 'semantic-ui-react';

const RecipeEditorInstructionItem = ({
  sectionIndex,
  instructionIndex,
  instruction,
  editInstruction,
  addInstruction,
  deleteInstruction,
  getCumulativeInstructionNum,
  moveInstruction
}) => (
  <Grid padded={false}>
    <Grid.Column width={2} className="shortGridColumn">
      <Label horizontal>
        {getCumulativeInstructionNum(sectionIndex, instructionIndex)}
      </Label>
    </Grid.Column>
    <Grid.Column width={8} className="shortGridColumn">
      <Input
        fluid
        type="text"
        value={instruction.name}
        onChange={e => editInstruction(sectionIndex, instructionIndex, e)}
      />
    </Grid.Column>
    <Grid.Column className="shortGridColumn">
      <Button.Group>
        <Button
          icon="arrow up"
          onClick={() => addInstruction(sectionIndex, instructionIndex)}
          content="New"
        />
        <Button
          icon="arrow down"
          onClick={() => addInstruction(sectionIndex, instructionIndex + 1)}
          content="New"
        />
        <Button
          icon="trash"
          onClick={() => deleteInstruction(sectionIndex, instructionIndex)}
        />
        <Button
          icon="arrow up"
          onClick={() =>
            moveInstruction(sectionIndex, instructionIndex, instructionIndex)
          }
        />
        <Button
          icon="arrow down"
          onClick={() =>
            moveInstruction(
              sectionIndex,
              instructionIndex,
              instructionIndex + 11
            )
          }
        />
      </Button.Group>
    </Grid.Column>
  </Grid>
);
export default RecipeEditorInstructionItem;

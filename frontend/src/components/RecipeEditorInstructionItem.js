import React from 'react';
import { Button, Grid, Input, Label } from 'semantic-ui-react';

const RecipeEditorInstructionItem = ({
  sectionNum,
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
        {getCumulativeInstructionNum(sectionNum, instructionIndex)}
      </Label>
    </Grid.Column>
    <Grid.Column width={8} className="shortGridColumn">
      <Input
        fluid
        type="text"
        value={instruction.name}
        onChange={e => editInstruction(sectionNum, instructionIndex, e)}
      />
    </Grid.Column>
    <Grid.Column className="shortGridColumn">
      <Button.Group>
        <Button
          icon="arrow up"
          onClick={() => addInstruction(sectionNum, instructionIndex)}
          content="New"
        />
        <Button
          icon="arrow down"
          onClick={() => addInstruction(sectionNum, instructionIndex + 1)}
          content="New"
        />
        <Button
          icon="trash"
          onClick={() => deleteInstruction(sectionNum, instructionIndex)}
        />
        <Button
          icon="arrow up"
          onClick={() =>
            moveInstruction(sectionNum, instructionIndex, instructionIndex)
          }
        />
        <Button
          icon="arrow down"
          onClick={() =>
            moveInstruction(sectionNum, instructionIndex, instructionIndex + 11)
          }
        />
      </Button.Group>
    </Grid.Column>
  </Grid>
);
export default RecipeEditorInstructionItem;

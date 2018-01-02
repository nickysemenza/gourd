import React from 'react';
import { Button, Grid, Input, Label } from 'semantic-ui-react';

const RecipeEditorInstructionItem = ({
  sectionNum,
  instructionNum,
  instruction,
  editInstruction,
  addInstruction,
  deleteInstruction,
  getCumulativeInstructionNum
}) => (
  <Grid padded={false}>
    <Grid.Column width={2} className="shortGridColumn">
      <Label horizontal>
        {getCumulativeInstructionNum(sectionNum, instructionNum)}
      </Label>
    </Grid.Column>
    <Grid.Column width={8} className="shortGridColumn">
      <Input
        fluid
        type="text"
        value={instruction.name}
        onChange={e => editInstruction(sectionNum, instructionNum, e)}
      />
    </Grid.Column>
    <Grid.Column className="shortGridColumn">
      <Button.Group>
        <Button
          icon="arrow up"
          onClick={() => addInstruction(sectionNum, instructionNum)}
          content="New"
        />
        <Button
          icon="arrow down"
          onClick={() => addInstruction(sectionNum, instructionNum + 1)}
          content="New"
        />
        <Button
          icon="trash"
          onClick={() => deleteInstruction(sectionNum, instructionNum)}
        />
      </Button.Group>
    </Grid.Column>
  </Grid>
);
export default RecipeEditorInstructionItem;

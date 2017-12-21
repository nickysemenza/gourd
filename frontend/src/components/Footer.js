import React from 'react';
import {
  Container,
  Grid,
  Header,
  Icon,
  List,
  Segment
} from 'semantic-ui-react';

const Nav = () => (
  <Segment inverted vertical style={{ padding: '5em 0em' }}>
    <Container>
      <Grid divided inverted stackable>
        <Grid.Row>
          <Grid.Column width={3}>
            <Header inverted as="h4" content="Info" />
            <List link inverted>
              <List.Item
                as="a"
                target="_blank"
                rel="noopener noreferrer"
                href="https://github.com/nickysemenza/food"
              >
                <Icon name="github" />Code on GitHub
              </List.Item>
              <List.Item
                as="a"
                target="_blank"
                rel="noopener noreferrer"
                href="https://www.nicky.photos/Food/My-Food/"
              >
                <Icon name="camera" />All Photos
              </List.Item>
              <List.Item
                as="a"
                target="_blank"
                rel="noopener noreferrer"
                href="https://instagram.com/nickysemenza"
              >
                <Icon name="instagram" />Instagram
              </List.Item>
              <List.Item as="a" href="mailto:nicky@nickysemenza.com">
                <Icon name="mail" />Contact
              </List.Item>
            </List>
          </Grid.Column>
          <Grid.Column width={7}>
            <Header as="h4" inverted>
              About
            </Header>
            <p>
              This is where I have my go-to recipes, ranging from baking to
              savory cooking. <br />Some are from various books and websites,
              others are my own, or somewhere in between.
            </p>
          </Grid.Column>
        </Grid.Row>
      </Grid>
    </Container>
  </Segment>
);
export default Nav;

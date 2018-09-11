import React, { Component } from 'react';

import { fetchImages } from '../actions/recipe';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { Header, Image, List, Table } from 'semantic-ui-react';
import { Link } from 'react-router-dom';

class ImageList extends Component {
  componentDidMount() {
    this.props.fetchImages();
  }
  render() {
    return (
      <div className="container">
        <Header dividing content="Images" />
        <Table celled>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Image</Table.HeaderCell>
              <Table.HeaderCell>dump</Table.HeaderCell>
              <Table.HeaderCell>Recipes</Table.HeaderCell>
            </Table.Row>
          </Table.Header>

          <Table.Body>
            {this.props.image_list.map(image => (
              <Table.Row key={image.id}>
                <Table.Cell>
                  <Image size="small" src={image.url} />
                </Table.Cell>
                <Table.Cell>
                  <pre>{JSON.stringify(image, true, 2)}</pre>
                </Table.Cell>
                <Table.Cell>
                  <List link>
                    {image.recipes.map(r => (
                      <List.Item key={r.id} as={Link} to={`/${r.slug}`}>
                        {r.title}
                      </List.Item>
                    ))}
                  </List>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </div>
    );
  }
}

function mapStateToProps(state) {
  let { image_list } = state.recipe;
  return { image_list };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      fetchImages
    },
    dispatch
  );
};

export default connect(mapStateToProps, mapDispatchToProps)(ImageList);

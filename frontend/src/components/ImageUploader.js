import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Dropzone from 'react-dropzone';
import { Button, Segment } from 'semantic-ui-react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { fetchRecipes } from '../actions/recipe';
import { API_BASE_URL } from '../config';
class ImageUploader extends Component {
  onResumeDrop(accepted, rejected) {
    if (rejected.length > 0) {
      console.log('wrong filetype for one or more files?');
      return;
    }
    console.log(accepted);

    let xhr = new XMLHttpRequest();

    xhr.onload = () => {
      // let { status } = xhr;
    };

    xhr.upload.addEventListener('progress', e => {
      console.log(e);
    });

    let fd = new FormData();
    for (let x = 0; x < accepted.length; x++) {
      fd.append('file', accepted[x]);
    }
    fd.append('slug', this.props.slug);
    xhr.open('PUT', API_BASE_URL + '/imageupload');
    xhr.send(fd);
  }
  render() {
    let dropzoneRef;
    return (
      <Segment>
        <Dropzone
          ref={node => {
            dropzoneRef = node;
          }}
          disableClick
          multiple={true}
          accept="image/*"
          onDrop={this.onResumeDrop.bind(this)}
          style={{ border: 'none', height: '100%' }}
        >
          hello
          <Button
            onClick={() => {
              dropzoneRef.open();
            }}
          >
            Drop or click to upload
          </Button>
        </Dropzone>
      </Segment>
    );
  }
}

function mapStateToProps(state) {
  return {
    recipe_list: state.recipe.recipe_list
  };
}

const mapDispatchToProps = dispatch => {
  return bindActionCreators(
    {
      fetchRecipes
    },
    dispatch
  );
};
ImageUploader.propTypes = {
  slug: PropTypes.string.isRequired
};
export default connect(mapStateToProps, mapDispatchToProps)(ImageUploader);

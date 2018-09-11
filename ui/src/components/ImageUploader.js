import React, { Component } from 'react';
import PropTypes from 'prop-types';
import Dropzone from 'react-dropzone';
import { Button, Progress, Segment } from 'semantic-ui-react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { fetchRecipes } from '../actions/recipe';
import { API_BASE_URL } from '../config';
import { toastr } from 'react-redux-toastr';

class ImageUploader extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isUploading: false,
      uploadProgress: 0
    };
  }
  onResumeDrop(accepted, rejected) {
    if (rejected.length > 0) {
      console.log('wrong filetype for one or more files?');
      return;
    }
    console.log(accepted);
    this.setState({ isUploading: true, uploadProgress: 0 });
    let xhr = new XMLHttpRequest();

    xhr.onload = () => {
      if (xhr.status === 200) {
        toastr.success('Success!', `${accepted.length} uploaded.`);
        this.props.onSuccessfulUpload();
      }
      setTimeout(() => this.setState({ isUploading: false }), 500);
    };

    xhr.upload.addEventListener('progress', event => {
      this.setState({ uploadProgress: event.loaded / event.total * 100 });
    });

    let fd = new FormData();
    for (let x = 0; x < accepted.length; x++) {
      fd.append('file', accepted[x]);
    }
    fd.append('slug', this.props.slug);
    xhr.open('PUT', API_BASE_URL + '/imageupload');
    xhr.setRequestHeader('X-JWT', this.props.jwt);
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
          {this.state.isUploading ? (
            <Progress percent={this.state.uploadProgress} />
          ) : null}
        </Dropzone>
      </Segment>
    );
  }
}

function mapStateToProps(state) {
  return {
    recipe_list: state.recipe.recipe_list,
    jwt: state.user.token
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
  slug: PropTypes.string.isRequired,
  onSuccessfulUpload: PropTypes.func
};
export default connect(mapStateToProps, mapDispatchToProps)(ImageUploader);

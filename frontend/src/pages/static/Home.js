import React, { Component } from 'react';
import RecipeList from '../../components/RecipeList';
import MainPhotoGallery from '../../components/MainPhotoGallery';
import { Divider, Header } from 'semantic-ui-react';

const Home = () => (
  <div>
    <Header dividing as="h1">
      Nicky's Recipe Stash
    </Header>
    <RecipeList />
    <Divider />
    <MainPhotoGallery />
  </div>
);
export default Home;

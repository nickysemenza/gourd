import gql from 'graphql-tag';
import * as React from 'react';
import * as ApolloReactCommon from '@apollo/react-common';
import * as ApolloReactComponents from '@apollo/react-components';
import * as ApolloReactHoc from '@apollo/react-hoc';
import * as ApolloReactHooks from '@apollo/react-hooks';
export type Maybe<T> = T | null;
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type Ingredient = {
   __typename?: 'Ingredient';
  uuid: Scalars['String'];
  name: Scalars['String'];
};

export type SectionInstruction = {
   __typename?: 'SectionInstruction';
  uuid: Scalars['String'];
  instruction: Scalars['String'];
};

export type SectionIngredient = {
   __typename?: 'SectionIngredient';
  uuid: Scalars['String'];
  info: Ingredient;
  grams: Scalars['Float'];
};

export type Section = {
   __typename?: 'Section';
  uuid: Scalars['String'];
  minutes: Scalars['Int'];
  instructions: Array<Maybe<SectionInstruction>>;
  ingredients: Array<Maybe<SectionIngredient>>;
};

export type Recipe = {
   __typename?: 'Recipe';
  uuid: Scalars['String'];
  name: Scalars['String'];
  total_minutes: Scalars['Int'];
  unit: Scalars['String'];
  sections: Array<Maybe<Section>>;
};

export type NewRecipe = {
  name: Scalars['String'];
};

export type Mutation = {
   __typename?: 'Mutation';
  createRecipe: Recipe;
};


export type MutationCreateRecipeArgs = {
  input?: Maybe<NewRecipe>;
};

export type Query = {
   __typename?: 'Query';
  recipes: Array<Recipe>;
  recipe?: Maybe<Recipe>;
};


export type QueryRecipeArgs = {
  uuid: Scalars['String'];
};

export type GetRecipeByUuidQueryVariables = {
  uuid: Scalars['String'];
};


export type GetRecipeByUuidQuery = (
  { __typename?: 'Query' }
  & { recipe?: Maybe<(
    { __typename?: 'Recipe' }
    & Pick<Recipe, 'uuid' | 'name' | 'total_minutes' | 'unit'>
    & { sections: Array<Maybe<(
      { __typename?: 'Section' }
      & Pick<Section, 'minutes'>
      & { ingredients: Array<Maybe<(
        { __typename?: 'SectionIngredient' }
        & Pick<SectionIngredient, 'uuid' | 'grams'>
        & { info: (
          { __typename?: 'Ingredient' }
          & Pick<Ingredient, 'name'>
        ) }
      )>>, instructions: Array<Maybe<(
        { __typename?: 'SectionInstruction' }
        & Pick<SectionInstruction, 'instruction' | 'uuid'>
      )>> }
    )>> }
  )> }
);

export type GetRecipesQueryVariables = {};


export type GetRecipesQuery = (
  { __typename?: 'Query' }
  & { recipes: Array<(
    { __typename: 'Recipe' }
    & Pick<Recipe, 'uuid' | 'name' | 'total_minutes' | 'unit'>
  )> }
);


export const GetRecipeByUuidDocument = gql`
    query getRecipeByUUID($uuid: String!) {
  recipe(uuid: $uuid) {
    uuid
    name
    total_minutes
    unit
    sections {
      minutes
      ingredients {
        uuid
        info {
          name
        }
        grams
      }
      instructions {
        instruction
        uuid
      }
    }
  }
}
    `;
export type GetRecipeByUuidComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<GetRecipeByUuidQuery, GetRecipeByUuidQueryVariables>, 'query'> & ({ variables: GetRecipeByUuidQueryVariables; skip?: boolean; } | { skip: boolean; });

    export const GetRecipeByUuidComponent = (props: GetRecipeByUuidComponentProps) => (
      <ApolloReactComponents.Query<GetRecipeByUuidQuery, GetRecipeByUuidQueryVariables> query={GetRecipeByUuidDocument} {...props} />
    );
    
export type GetRecipeByUuidProps<TChildProps = {}, TDataName extends string = 'data'> = {
      [key in TDataName]: ApolloReactHoc.DataValue<GetRecipeByUuidQuery, GetRecipeByUuidQueryVariables>
    } & TChildProps;
export function withGetRecipeByUuid<TProps, TChildProps = {}, TDataName extends string = 'data'>(operationOptions?: ApolloReactHoc.OperationOption<
  TProps,
  GetRecipeByUuidQuery,
  GetRecipeByUuidQueryVariables,
  GetRecipeByUuidProps<TChildProps, TDataName>>) {
    return ApolloReactHoc.withQuery<TProps, GetRecipeByUuidQuery, GetRecipeByUuidQueryVariables, GetRecipeByUuidProps<TChildProps, TDataName>>(GetRecipeByUuidDocument, {
      alias: 'getRecipeByUuid',
      ...operationOptions
    });
};

/**
 * __useGetRecipeByUuidQuery__
 *
 * To run a query within a React component, call `useGetRecipeByUuidQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetRecipeByUuidQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetRecipeByUuidQuery({
 *   variables: {
 *      uuid: // value for 'uuid'
 *   },
 * });
 */
export function useGetRecipeByUuidQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<GetRecipeByUuidQuery, GetRecipeByUuidQueryVariables>) {
        return ApolloReactHooks.useQuery<GetRecipeByUuidQuery, GetRecipeByUuidQueryVariables>(GetRecipeByUuidDocument, baseOptions);
      }
export function useGetRecipeByUuidLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetRecipeByUuidQuery, GetRecipeByUuidQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<GetRecipeByUuidQuery, GetRecipeByUuidQueryVariables>(GetRecipeByUuidDocument, baseOptions);
        }
export type GetRecipeByUuidQueryHookResult = ReturnType<typeof useGetRecipeByUuidQuery>;
export type GetRecipeByUuidLazyQueryHookResult = ReturnType<typeof useGetRecipeByUuidLazyQuery>;
export type GetRecipeByUuidQueryResult = ApolloReactCommon.QueryResult<GetRecipeByUuidQuery, GetRecipeByUuidQueryVariables>;
export const GetRecipesDocument = gql`
    query getRecipes {
  recipes {
    uuid
    name
    total_minutes
    unit
    __typename
  }
}
    `;
export type GetRecipesComponentProps = Omit<ApolloReactComponents.QueryComponentOptions<GetRecipesQuery, GetRecipesQueryVariables>, 'query'>;

    export const GetRecipesComponent = (props: GetRecipesComponentProps) => (
      <ApolloReactComponents.Query<GetRecipesQuery, GetRecipesQueryVariables> query={GetRecipesDocument} {...props} />
    );
    
export type GetRecipesProps<TChildProps = {}, TDataName extends string = 'data'> = {
      [key in TDataName]: ApolloReactHoc.DataValue<GetRecipesQuery, GetRecipesQueryVariables>
    } & TChildProps;
export function withGetRecipes<TProps, TChildProps = {}, TDataName extends string = 'data'>(operationOptions?: ApolloReactHoc.OperationOption<
  TProps,
  GetRecipesQuery,
  GetRecipesQueryVariables,
  GetRecipesProps<TChildProps, TDataName>>) {
    return ApolloReactHoc.withQuery<TProps, GetRecipesQuery, GetRecipesQueryVariables, GetRecipesProps<TChildProps, TDataName>>(GetRecipesDocument, {
      alias: 'getRecipes',
      ...operationOptions
    });
};

/**
 * __useGetRecipesQuery__
 *
 * To run a query within a React component, call `useGetRecipesQuery` and pass it any options that fit your needs.
 * When your component renders, `useGetRecipesQuery` returns an object from Apollo Client that contains loading, error, and data properties
 * you can use to render your UI.
 *
 * @param baseOptions options that will be passed into the query, supported options are listed on: https://www.apollographql.com/docs/react/api/react-hooks/#options;
 *
 * @example
 * const { data, loading, error } = useGetRecipesQuery({
 *   variables: {
 *   },
 * });
 */
export function useGetRecipesQuery(baseOptions?: ApolloReactHooks.QueryHookOptions<GetRecipesQuery, GetRecipesQueryVariables>) {
        return ApolloReactHooks.useQuery<GetRecipesQuery, GetRecipesQueryVariables>(GetRecipesDocument, baseOptions);
      }
export function useGetRecipesLazyQuery(baseOptions?: ApolloReactHooks.LazyQueryHookOptions<GetRecipesQuery, GetRecipesQueryVariables>) {
          return ApolloReactHooks.useLazyQuery<GetRecipesQuery, GetRecipesQueryVariables>(GetRecipesDocument, baseOptions);
        }
export type GetRecipesQueryHookResult = ReturnType<typeof useGetRecipesQuery>;
export type GetRecipesLazyQueryHookResult = ReturnType<typeof useGetRecipesLazyQuery>;
export type GetRecipesQueryResult = ApolloReactCommon.QueryResult<GetRecipesQuery, GetRecipesQueryVariables>;
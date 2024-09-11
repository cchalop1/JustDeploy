/* eslint-disable */
import { DocumentTypeDecoration } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
};

export type App = {
  __typename?: 'App';
  id: Scalars['ID']['output'];
  isDockerCompose: Scalars['Boolean']['output'];
  isDockerFile: Scalars['Boolean']['output'];
  name: Scalars['String']['output'];
  path: Scalars['String']['output'];
};

export type Env = {
  __typename?: 'Env';
  name: Scalars['String']['output'];
  value: Scalars['String']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createService: Service;
};


export type MutationCreateServiceArgs = {
  input: NewService;
};

export type NewService = {
  deployId?: InputMaybe<Scalars['String']['input']>;
  fromDockerCompose: Scalars['Boolean']['input'];
  projectId?: InputMaybe<Scalars['String']['input']>;
  serviceName: Scalars['String']['input'];
};

export type Project = {
  __typename?: 'Project';
  apps: Array<App>;
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  path: Scalars['String']['output'];
  services: Array<Service>;
};

export type Query = {
  __typename?: 'Query';
  getProject?: Maybe<Project>;
};


export type QueryGetProjectArgs = {
  id: Scalars['ID']['input'];
};

export type Service = {
  __typename?: 'Service';
  deployId?: Maybe<Scalars['String']['output']>;
  envs: Array<Env>;
  host: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  imageName: Scalars['String']['output'];
  imageUrl: Scalars['String']['output'];
  name: Scalars['String']['output'];
  projectId?: Maybe<Scalars['String']['output']>;
  status: Scalars['String']['output'];
  volumsNames: Array<Scalars['String']['output']>;
};

export type GetProjectQueryQueryVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type GetProjectQueryQuery = { __typename?: 'Query', getProject?: { __typename?: 'Project', id: string, name: string, path: string, apps: Array<{ __typename?: 'App', id: string, name: string, path: string }>, services: Array<{ __typename?: 'Service', id: string, name: string, status: string }> } | null };

export class TypedDocumentString<TResult, TVariables>
  extends String
  implements DocumentTypeDecoration<TResult, TVariables>
{
  __apiType?: DocumentTypeDecoration<TResult, TVariables>['__apiType'];

  constructor(private value: string, public __meta__?: Record<string, any>) {
    super(value);
  }

  toString(): string & DocumentTypeDecoration<TResult, TVariables> {
    return this.value;
  }
}

export const GetProjectQueryDocument = new TypedDocumentString(`
    query getProjectQuery($id: ID!) {
  getProject(id: $id) {
    id
    name
    path
    apps {
      id
      name
      path
    }
    services {
      id
      name
      status
    }
  }
}
    `) as unknown as TypedDocumentString<GetProjectQueryQuery, GetProjectQueryQueryVariables>;
/* eslint-disable */
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type Cart = {
  __typename?: 'Cart';
  id: Scalars['String'];
};

export type CartItem = {
  productId: Scalars['String'];
};

export type Game = {
  __typename?: 'Game';
  category: Scalars['String'];
  id: Scalars['String'];
  image_url: Scalars['String'];
  price: Scalars['String'];
  product: Scalars['String'];
};

export type Inventory = {
  __typename?: 'Inventory';
  games: Array<Game>;
};

export type Mutation = {
  __typename?: 'Mutation';
  addGameToCart: Cart;
};


export type MutationAddGameToCartArgs = {
  input: CartItem;
};

export type Query = {
  __typename?: 'Query';
  inventory: Inventory;
  shopper: Shopper;
};


export type QueryShopperArgs = {
  input?: InputMaybe<ShopperInput>;
};

export type Shopper = {
  __typename?: 'Shopper';
  email: Scalars['String'];
  id: Scalars['String'];
  inventoryId: Scalars['String'];
};

export type ShopperInput = {
  shopperId?: InputMaybe<Scalars['String']>;
};

/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type Cart = {
  __typename?: 'Cart';
  id: Scalars['String'];
};

export type CartItem = {
  productId: Scalars['String'];
};

export type Game = {
  __typename?: 'Game';
  category: Scalars['String'];
  id: Scalars['String'];
  image_url: Scalars['String'];
  price: Scalars['String'];
  product: Scalars['String'];
};

export type Inventory = {
  __typename?: 'Inventory';
  games: Array<Game>;
};

export type Mutation = {
  __typename?: 'Mutation';
  addGameToCart: Cart;
};


export type MutationAddGameToCartArgs = {
  input: CartItem;
};

export type Query = {
  __typename?: 'Query';
  inventory: Inventory;
  shopper: Shopper;
};


export type QueryShopperArgs = {
  input?: InputMaybe<ShopperInput>;
};

export type Shopper = {
  __typename?: 'Shopper';
  email: Scalars['String'];
  id: Scalars['String'];
  inventoryId: Scalars['String'];
};

export type ShopperInput = {
  shopperId?: InputMaybe<Scalars['String']>;
};

import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
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
  Time: any;
};

export type Cart = {
  __typename?: 'Cart';
  id: Scalars['String'];
  items: Array<CartItem>;
  shopperId: Scalars['String'];
  subtotal?: Maybe<Scalars['String']>;
  tax?: Maybe<Scalars['String']>;
  taxRate?: Maybe<Scalars['String']>;
  timestamp?: Maybe<Scalars['Time']>;
  total?: Maybe<Scalars['String']>;
};

export type CartInput = {
  cartId: Scalars['String'];
  shopperId: Scalars['String'];
};

export type CartItem = {
  __typename?: 'CartItem';
  category: Scalars['String'];
  price: Scalars['String'];
  productId: Scalars['String'];
  quantity: Scalars['Int'];
  subtotal: Scalars['String'];
  title: Scalars['String'];
};

export type CartItemInput = {
  productId: Scalars['String'];
  quantity: Scalars['Int'];
};

export type CartSubscriptionInput = {
  cartId: Scalars['String'];
  topic?: InputMaybe<Scalars['String']>;
};

export type Game = {
  __typename?: 'Game';
  category: Scalars['String'];
  id: Scalars['String'];
  imageUrl: Scalars['String'];
  price: Scalars['String'];
  title: Scalars['String'];
};

export type Inventory = {
  __typename?: 'Inventory';
  categories: Array<Scalars['String']>;
  games: Array<Game>;
};

export type InventoryInput = {
  category?: InputMaybe<Scalars['String']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  setCartItems: Cart;
};


export type MutationSetCartItemsArgs = {
  input?: InputMaybe<SetCartItemsInput>;
};

export type Ping = {
  __typename?: 'Ping';
  timestamp?: Maybe<Scalars['Time']>;
  value?: Maybe<Scalars['String']>;
};

export type PingInput = {
  timestamp?: InputMaybe<Scalars['Time']>;
  value?: InputMaybe<Scalars['String']>;
};

export type Pong = {
  __typename?: 'Pong';
  timestamp?: Maybe<Scalars['Time']>;
  value?: Maybe<Scalars['String']>;
};

export type Query = {
  __typename?: 'Query';
  cart: Cart;
  inventory: Inventory;
  ping: Pong;
  shopper: Shopper;
  user: User;
};


export type QueryCartArgs = {
  input?: InputMaybe<CartInput>;
};


export type QueryInventoryArgs = {
  input?: InputMaybe<InventoryInput>;
};


export type QueryPingArgs = {
  input?: InputMaybe<PingInput>;
};


export type QueryShopperArgs = {
  input?: InputMaybe<ShopperInput>;
};


export type QueryUserArgs = {
  input?: InputMaybe<UserInput>;
};

export type SetCartItemsInput = {
  cartId: Scalars['String'];
  items: Array<CartItemInput>;
};

export type Shopper = {
  __typename?: 'Shopper';
  cartId: Scalars['String'];
  email: Scalars['String'];
  id: Scalars['String'];
  inventoryId: Scalars['String'];
};

export type ShopperInput = {
  shopperId?: InputMaybe<Scalars['String']>;
};

export type Subscription = {
  __typename?: 'Subscription';
  cart: Cart;
};


export type SubscriptionCartArgs = {
  input: CartSubscriptionInput;
};

export type User = {
  __typename?: 'User';
  email: Scalars['String'];
  ok: Scalars['Boolean'];
  token?: Maybe<Scalars['String']>;
};

export type UserInput = {
  token?: InputMaybe<Scalars['String']>;
};

export type SetCartItemsMutationVariables = Exact<{
  input: SetCartItemsInput;
}>;


export type SetCartItemsMutation = { __typename?: 'Mutation', setCartItems: { __typename?: 'Cart', id: string, shopperId: string, subtotal?: string | null, total?: string | null, taxRate?: string | null, tax?: string | null, timestamp?: any | null, items: Array<{ __typename?: 'CartItem', productId: string, quantity: number, price: string, subtotal: string, title: string, category: string }> } };

export type CurrentCartQueryVariables = Exact<{
  input?: InputMaybe<CartInput>;
}>;


export type CurrentCartQuery = { __typename?: 'Query', cart: { __typename?: 'Cart', id: string, shopperId: string, total?: string | null, subtotal?: string | null, taxRate?: string | null, tax?: string | null, timestamp?: any | null, items: Array<{ __typename?: 'CartItem', productId: string, quantity: number, price: string, subtotal: string, title: string, category: string }> } };

export type CurrentUserQueryVariables = Exact<{
  input?: InputMaybe<UserInput>;
}>;


export type CurrentUserQuery = { __typename?: 'Query', user: { __typename?: 'User', email: string, token?: string | null, ok: boolean } };

export type InventoryQueryVariables = Exact<{
  input: InventoryInput;
}>;


export type InventoryQuery = { __typename?: 'Query', inventory: { __typename?: 'Inventory', categories: Array<string>, games: Array<{ __typename?: 'Game', id: string, title: string, imageUrl: string, category: string, price: string }> } };

export type PingTestQueryVariables = Exact<{
  input?: InputMaybe<PingInput>;
}>;


export type PingTestQuery = { __typename?: 'Query', ping: { __typename?: 'Pong', value?: string | null, timestamp?: any | null } };

export type CartSubscriptionVariables = Exact<{
  input: CartSubscriptionInput;
}>;


export type CartSubscription = { __typename?: 'Subscription', cart: { __typename?: 'Cart', id: string, shopperId: string, total?: string | null, subtotal?: string | null, taxRate?: string | null, tax?: string | null, timestamp?: any | null, items: Array<{ __typename?: 'CartItem', productId: string, quantity: number, price: string, subtotal: string, title: string, category: string }> } };


export const SetCartItemsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"SetCartItems"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"SetCartItemsInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"setCartItems"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"shopperId"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"productId"}},{"kind":"Field","name":{"kind":"Name","value":"quantity"}},{"kind":"Field","name":{"kind":"Name","value":"price"}},{"kind":"Field","name":{"kind":"Name","value":"subtotal"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"category"}}]}},{"kind":"Field","name":{"kind":"Name","value":"subtotal"}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"taxRate"}},{"kind":"Field","name":{"kind":"Name","value":"tax"}},{"kind":"Field","name":{"kind":"Name","value":"timestamp"}}]}}]}}]} as unknown as DocumentNode<SetCartItemsMutation, SetCartItemsMutationVariables>;
export const CurrentCartDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"CurrentCart"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"CartInput"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"cart"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"shopperId"}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"subtotal"}},{"kind":"Field","name":{"kind":"Name","value":"taxRate"}},{"kind":"Field","name":{"kind":"Name","value":"tax"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"productId"}},{"kind":"Field","name":{"kind":"Name","value":"quantity"}},{"kind":"Field","name":{"kind":"Name","value":"price"}},{"kind":"Field","name":{"kind":"Name","value":"subtotal"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"category"}}]}},{"kind":"Field","name":{"kind":"Name","value":"timestamp"}}]}}]}}]} as unknown as DocumentNode<CurrentCartQuery, CurrentCartQueryVariables>;
export const CurrentUserDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"CurrentUser"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"UserInput"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"user"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"email"}},{"kind":"Field","name":{"kind":"Name","value":"token"}},{"kind":"Field","name":{"kind":"Name","value":"ok"}}]}}]}}]} as unknown as DocumentNode<CurrentUserQuery, CurrentUserQueryVariables>;
export const InventoryDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"Inventory"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"InventoryInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"inventory"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"categories"}},{"kind":"Field","name":{"kind":"Name","value":"games"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"imageUrl"}},{"kind":"Field","name":{"kind":"Name","value":"category"}},{"kind":"Field","name":{"kind":"Name","value":"price"}}]}}]}}]}}]} as unknown as DocumentNode<InventoryQuery, InventoryQueryVariables>;
export const PingTestDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"PingTest"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"PingInput"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"ping"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"value"}},{"kind":"Field","name":{"kind":"Name","value":"timestamp"}}]}}]}}]} as unknown as DocumentNode<PingTestQuery, PingTestQueryVariables>;
export const CartDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"Cart"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"input"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CartSubscriptionInput"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"cart"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"input"},"value":{"kind":"Variable","name":{"kind":"Name","value":"input"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"shopperId"}},{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"subtotal"}},{"kind":"Field","name":{"kind":"Name","value":"taxRate"}},{"kind":"Field","name":{"kind":"Name","value":"tax"}},{"kind":"Field","name":{"kind":"Name","value":"items"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"productId"}},{"kind":"Field","name":{"kind":"Name","value":"quantity"}},{"kind":"Field","name":{"kind":"Name","value":"price"}},{"kind":"Field","name":{"kind":"Name","value":"subtotal"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"category"}}]}},{"kind":"Field","name":{"kind":"Name","value":"timestamp"}}]}}]}}]} as unknown as DocumentNode<CartSubscription, CartSubscriptionVariables>;
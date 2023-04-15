import type {
    AnyVariables,
    Client,
    DocumentInput,
    GraphQLRequest,
    Operation,
    OperationContext, OperationResult, OperationResultSource,
    OperationType
} from "@urql/svelte";
import type { Sink} from 'wonka'
import {fromValue, never} from "wonka";

let mc


// mc = {
//     executeQuery: ({query}) => {
//         console.log('BONKBONKBONK')
//         if(query == 'CurrentCart') {
//             return fromValue(initialCart)
//         }
//         throw new Error('oops')
//     },
//     executeMutation: vi.fn(() => never),
//     executeSubscription: vi.fn(() => never),
//     operations$: vi.fn(()=>never),
//     suspense: false,
//     reexecuteOperation: vi.fn(()=>never),
//     createRequestOperation: vi.fn(()=>never),
//     executeRequestOperation: vi.fn(() => never),
//     query: ({query}) => {
//         console.log('BONKBONKBONK')
//         if(query == 'CurrentCart') {
//             return fromValue(initialCart)
//         }
//         throw new Error('oops')
//     },
//     mutation: vi.fn(()=>never),
//     readQuery: vi.fn(()=>never),
//     subscription: vi.fn(()=>never),
// }
    // operations$(sink: Sink<Operation<any, AnyVariables>>): void {
    // },
    // suspense: false,
    // createRequestOperation<Data = any, Variables = AnyVariables>(kind: OperationType, request: GraphQLRequest<Data, Variables>, opts?: Partial<OperationContext> | undefined): Operation<Data, Variables> {
    //     return undefined;
    // },
    // executeMutation<Data = any, Variables = AnyVariables extends AnyVariables>(query: GraphQLRequest<Data, Variables>, opts?: Partial<OperationContext> | undefined): OperationResultSource<OperationResult<Data, Variables>> {
    //     return undefined;
    // },
    // executeQuery<Data = any, Variables = AnyVariables extends AnyVariables>(query: GraphQLRequest<Data, Variables>, opts?: Partial<OperationContext> | undefined): OperationResultSource<OperationResult<Data, Variables>> {
    //     return undefined;
    // },
    // executeRequestOperation<Data = any, Variables = AnyVariables extends AnyVariables>(operation: Operation<Data, Variables>): OperationResultSource<OperationResult<Data, Variables>> {
    //     return undefined;
    // },
    // executeSubscription<Data = any, Variables = AnyVariables extends AnyVariables>(query: GraphQLRequest<Data, Variables>, opts?: Partial<OperationContext> | undefined): OperationResultSource<OperationResult<Data, Variables>> {
    //     return undefined;
    // },
    // mutation<Data = any, Variables = AnyVariables extends AnyVariables>(query: DocumentInput<Data, Variables>, variables: Variables, context?: Partial<OperationContext>): OperationResultSource<OperationResult<Data, Variables>> {
    //     return undefined;
    // },
    // query<Data = any, Variables = AnyVariables extends AnyVariables>(query: DocumentInput<Data, Variables>, variables: Variables, context?: Partial<OperationContext>): OperationResultSource<OperationResult<Data, Variables>> {
    //     return undefined;
    // },
    // readQuery<Data = any, Variables = AnyVariables extends AnyVariables>(query: DocumentInput<Data, Variables>, variables: Variables, context?: Partial<OperationContext>): OperationResult<Data, Variables> | null {
    //     return undefined;
    // },
    // reexecuteOperation(operation: Operation): void {
    // },
    // subscription<Data = any, Variables = AnyVariables extends AnyVariables>(query: DocumentInput<Data, Variables>, variables: Variables, context?: Partial<OperationContext>): OperationResultSource<OperationResult<Data, Variables>> {
    //     return undefined;
    // }
}
import { stringify } from 'yaml'
import type {CodegenConfig} from '@graphql-codegen/cli'
import {writeFileSync} from 'fs'

const config: CodegenConfig = {
    schema: '../../graphql/schema.graphqls',
    documents: ['src/**/*.svelte'],
    ignoreNoDocuments: true,
    generates: {
        './src/gql/': {
            preset: 'client-preset',
            plugins: ['typescript','typescript-operations','typed-document-node'],
            config: {
                useTypeImports: true
            }
        }
    },
}

//
// https://github.com/dotansimha/graphql-code-generator/issues/8488#issuecomment-1340622934
// save config as yml since TS5095 warning will be raised if using codegen.ts directly
writeFileSync('codegen.yml', stringify(config));

export default config
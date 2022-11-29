<script>
    import { PUBLIC_GRAPHQL_URL} from "$env/static/public";
    import { queryStore, gql, getContextClient } from '@urql/svelte';

    const todos = queryStore({
        client: getContextClient(),
        query: gql`
      query {
        todos {
          id
          title
        }
      }
    `,
    });
</script>
<h1>Temporal Shop</h1>


<p>Connecting to graphql at {PUBLIC_GRAPHQL_URL}</p>
<h3>Routes</h3>
<ul>
    <li><a href="/foo">foo</a></li>
    <li><a href="/foo/deep-foo">deep foo</a></li>
</ul>
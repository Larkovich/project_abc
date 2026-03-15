<script>
  import Router from 'svelte-spa-router';
  import Dashboard from './Dashboard.svelte';
  import MagicLink from './MagicLink.svelte';

  const routes = {
    '/': Dashboard,
    '/p/:id': MagicLink,
  };
</script>

<Router {routes} />

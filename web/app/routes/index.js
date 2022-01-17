import Route from '@ember/routing/route';

export default class IndexRoute extends Route {
  queryParams = {
    page: {
      refreshModel: true,
      type: 'number',
    },
    limit: {
      refreshModel: true,
      type: 'number',
    },
  };
}

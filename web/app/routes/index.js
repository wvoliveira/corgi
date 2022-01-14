import Route from '@ember/routing/route';
import ENV from 'corgi/config/environment';

export default class IndexRoute extends Route {
  queryParams = {
    page: {
      refreshModel: true,
    },
  };
  async model(params) {
    let response = await fetch(`${ENV.APP.apiHost}/api/v1/links/`, {
      credentials: 'include',
    });
    let { data } = await response.json();

    let { id, attributes } = data;
    let type;

    return { id, type, ...attributes };
  }
}

import ENV from 'corgi/config/environment';
import RESTAdapter from '@ember-data/adapter/rest';

export default class ApplicationAdapter extends RESTAdapter {
  namespace = 'api/v1';

  buildURL(...args) {
    return `${ENV.APP.apiHost}${super.buildURL(...args)}`;
  }

  shouldReloadAll(store, snapshotsArray) {
    return false;
  }

  shouldBackgroundReloadAll(store, snapshotsArray) {
    return true;
  }
}

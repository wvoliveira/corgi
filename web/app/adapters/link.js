import ApplicationAdapter from './application';
import { inject as service } from '@ember/service';

export default class LinkAdapter extends ApplicationAdapter {
  @service session;

  findAll(store, type) {
    let url = this.buildURL('links');
    console.log('url: ' + url);

    return fetch(url, {
      method: 'GET',
      credentials: 'include',
    });
  }
}

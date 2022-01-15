import ApplicationAdapter from './application';

export default class LinkAdapter extends ApplicationAdapter {
  async fetch(url, method = 'GET') {
    let response = await fetch(url, { method: method, credentials: 'include' });
    console.log(response);
    return response.json();
  }

  findAll(store, type, sinceToken, snapshotRecordArray) {
    let url = this.buildURL('links');
    return this.fetch(url);
  }
}

import ApplicationAdapter from './application';

export default class LinkAdapter extends ApplicationAdapter {
  async fetch(url, method = 'GET', options) {
    let response = await fetch(url, { method: method, credentials: 'include', options });
    console.log(response);
    return response.json();
  }

  findAll(store, type, sinceToken, snapshotRecordArray) {
    let url = this.buildURL('links');
    return this.fetch(url);
  }

  query(store, type, query, recordArray, adapterOptions) {
    let URLparams = new URLSearchParams({
      page: query.page,
      limit: query.limit,
    });

    let url = this.buildURL('links?' + URLparams);
    return this.fetch(url);
  }
}

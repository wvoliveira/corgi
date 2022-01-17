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

  query(store, type, query, recordArray, adapterOptions) {
    console.log('adapter query');
    console.log('store');
    console.log(store);

    console.log('type');
    console.log(type);

    console.log('query');
    console.log(query);

    console.log('recordArray');
    console.log(recordArray);

    console.log('adapterOptions');
    console.log(adapterOptions);

    let url = this.buildURL('links');
    return this.fetch(url);
  }
}

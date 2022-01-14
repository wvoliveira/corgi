import ENV from 'corgi/config/environment';
import RESTAdapter from '@ember-data/adapter/rest';
import RSVP from 'rsvp';
import $ from 'jquery';

export default class LinkAdapter extends RESTAdapter {
  namespace = 'api/v1';

  buildURL(...args) {
    return `${ENV.APP.apiHost}${super.buildURL(...args)}/`;
  }

  findAll(store, type) {
    let url = 'http://localhost:8081/api/v1/links/';

    $.ajaxSetup({
      dataType: 'json',
      xhrFields: {
        withCredentials: true,
      },
      crossDomain: true,
    });

    return new RSVP.Promise(function(resolve, reject) {
      $.getJSON(url).then(function(data) {
        resolve(data);
      }, function(jqXHR) {
        reject(jqXHR);
      });
    });
  }
}

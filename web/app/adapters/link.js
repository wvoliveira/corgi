import RESTAdapter from '@ember-data/adapter/rest';
import RSVP from 'rsvp';
import $ from 'jquery';
import ENV from 'corgi/config/environment';

export default class LinkAdapter extends RESTAdapter {
  constructor() {
    super();
  }

  findAll(store, type) {
    console.log('find all');

    return new RSVP.Promise(function (resolve, reject) {
      // eslint-disable-next-line ember/no-jquery
      $.ajax({
        url: `${ENV.APP.apiHost}/api/v1/links/`,
        xhrFields: {
          withCredentials: true,
        },
        type: 'GET',
        processData: false,
        success: function (response) {
          console.log('response');
          console.log(response);
          resolve(response);
        },
        error: function (xhr, ajaxOptions, thrownError) {
          //Add these parameters to display the required response
          console.log(thrownError);
          reject(thrownError);
        },
      });

      // $.getJSON(`${ENV.APP.apiHost}/api/v1/links/`).then(
      //   function (data) {
      //     resolve(data);
      //   },
      //   function (jqXHR) {
      //     reject(jqXHR);
      //   }
      // );
    });
  }
}

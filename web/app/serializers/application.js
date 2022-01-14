import { underscore } from '@ember/string';
import RESTSerializer from '@ember-data/serializer/rest';
import ENV from 'corgi/config/environment';

export default class ApplicationSerializer extends RESTSerializer {
  apiEndpoint = ENV.APP.apiHost;

  keyForAttribute(attr) {
    return underscore(attr);
  }
}

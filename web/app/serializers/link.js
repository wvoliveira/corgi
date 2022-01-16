import RESTSerializer from '@ember-data/serializer/rest';

export default class LinkSerializer extends RESTSerializer {
  normalizeFindAllResponse(store, primaryModelClass, payload, id, requestType) {
    console.log('store: ' + store);
    console.log('primaryModelClass: ' + primaryModelClass);
    console.log('payload: ' + payload);
    console.log(payload);
    console.log('id: ' + id);
    console.log('requestType: ' + requestType);

    return super.normalizeFindAllResponse(...arguments);
  }
}

import RESTSerializer, {EmbeddedRecordsMixin} from '@ember-data/serializer/rest';

export default class LinkSerializer extends RESTSerializer.extend(
  EmbeddedRecordsMixin
) {
  normalizeFindAllResponse(store, primaryModelClass, payload, id, requestType) {
    let data = super.normalizeFindAllResponse(store, primaryModelClass, payload, id, requestType);
    console.log('data');
    console.log(data);

    console.log('store: ' + store);
    console.log('primaryModelClass: ' + primaryModelClass);
    console.log('payload: ' + payload);
    console.log(payload);
    console.log('id: ' + id);
    console.log('requestType: ' + requestType);

    return super.normalizeFindAllResponse(...arguments);
  }

  normalizeQueryResponse(store, primaryModelClass, payload, id, requestType) {
    let data = {
      links: payload.data,
      meta: {
        limit: payload.limit,
        page: payload.page,
        sort: payload.sort,
        total: payload.total,
        pages: payload.pages,
      },
    };

    payload = data;
    return super.normalizeQueryResponse(store, primaryModelClass, payload, id, requestType);
  }
}

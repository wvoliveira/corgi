import { underscore } from '@ember/string';
import RESTSerializer from '@ember-data/serializer/json-api';

export default class LinkSerializer extends RESTSerializer {
  primaryKey = 'id';
  keyForAttribute(attr) {
    return underscore(attr);
  }

  normalizeResponse(store, primaryModelClass, payload, id, requestType) {
    console.log('store');
    console.log(store);

    console.log('primaryModelClass');
    console.log(primaryModelClass);

    console.log('payload');
    console.log(payload);

    console.log('id');
    console.log(id);

    console.log('requestType');
    console.log(requestType);

    //payload.data.attributes.amount = payload.data.attributes.cost.amount;
    //payload.data.attributes.currency = payload.data.attributes.cost.currency;

    //return super.normalizeResponse(...arguments);

    //{"data":[{"id":"de071e6b-9eb3-4c31-af48-e76a528e00c2","created_at":"2021-12-08T17:06:44.88689-03:00","updated_at":"2021-12-08T17:06:44.88839-03:00","domain":"elga.io","keyword":"xablau","url":"https://elga.io","title":"ELGA home","active":"true"},{"id":"dab8ea12-7653-472d-bda1-82bc83f005b0","created_at":"2021-12-08T18:24:51.505437-03:00","updated_at":"2021-12-08T18:24:51.506437-03:00","domain":"elga.io","keyword":"abecedo112","url":"https://elga.io","title":"ELGA home","active":"true"},{"id":"c06d9c48-2fe8-4c93-ac3a-8c3abc47810f","created_at":"2022-01-08T14:59:32.582667-03:00","updated_at":"2022-01-08T14:59:32.583452-03:00","domain":"elga.io","keyword":"tutorialdobom","url":"https://medium.com/swlh/building-a-user-auth-system-with-jwt-using-golang-30892659cc0","title":"","active":"true"},{"id":"c0451a66-c877-4b3f-b76f-feb9c74d696f","created_at":"2022-01-10T21:42:06.120878-03:00","updated_at":"2022-01-10T21:42:06.121729-03:00","domain":"elga.io","keyword":"linkedincom","url":"https://www.linkedin.com/","title":"","active":"true"},{"id":"a0f232ed-60d2-4308-a005-52bf5a9c3977","created_at":"2021-12-08T18:24:54.100934-03:00","updated_at":"2021-12-08T18:24:54.101936-03:00","domain":"elga.io","keyword":"abecedo122312","url":"https://elga.io","title":"ELGA home","active":"true"},{"id":"9f34f06a-1b83-498c-9a20-fc59fe91b1f3","created_at":"2021-12-08T18:25:08.668908-03:00","updated_at":"2021-12-08T18:25:08.669907-03:00","domain":"elga.io","keyword":"axxxAASDXZZA2","url":"https://elga.io","title":"ELGA home","active":"true"},{"id":"9ab85263-5b6b-4cdf-a21e-e3000e7bce6e","created_at":"2021-12-08T18:46:45.303229-03:00","updated_at":"2021-12-08T18:46:45.30423-03:00","domain":"elga.io","keyword":"Asdferexz11","url":"https://elga.io","title":"ELGA home","active":"true"},{"id":"94bb7a90-4c3a-40c8-b7c5-2b4e5439cb9c","created_at":"2021-12-08T17:35:35.720646-03:00","updated_at":"2021-12-08T17:35:35.721647-03:00","domain":"elga.io","keyword":"abecedario","url":"https://elga.io","title":"ELGA home","active":"true"},{"id":"5ee6f54f-ccd7-4431-8bc1-8415bd94c530","created_at":"2021-12-08T18:51:43.81537-03:00","updated_at":"2021-12-08T18:51:43.816395-03:00","domain":"elga.io","keyword":"Asdferexz1000","url":"https://elga.io","title":"ELGA home","active":"true"},{"id":"4cce6c81-c9fb-49f1-ae8d-0235cd67d4af","created_at":"2021-12-08T18:24:57.440419-03:00","updated_at":"2021-12-08T18:24:57.44142-03:00","domain":"elga.io","keyword":"axxxxxx2","url":"https://elga.io","title":"ELGA home","active":"true"}]
    // "limit":10,"page":1,"sort":"ID desc","total":15,"pages":2}
    let d = {
      data: {
        id: 'links',
        type: 'link',
        data: payload.data,
        meta: {
          limit: payload.limit,
          page: payload.page,
          sort: payload.sort,
          total: payload.total,
          pages: payload.pages,
        },
      },
    };

    console.log('d:');
    console.log(d);

    return d;
  }
}

import Model, { attr } from '@ember-data/model';

export default class Pagination extends Model {
  @attr('number') Limit;
  @attr('number') Page;
  @attr('number') Pages;
  @attr('string') Sort;
  @attr('number') Total;
}

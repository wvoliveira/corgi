import Model, { attr } from '@ember-data/model';

export default class LinkModel extends Model {
  @attr data;
  @attr('number') limit;
  @attr('number') page;
  @attr('string') sort;
  @attr('number') total;
  @attr('number') pages;
}

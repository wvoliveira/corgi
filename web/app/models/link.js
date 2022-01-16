import Model, { attr } from '@ember-data/model';

export default class LinkModel extends Model {
  @attr data;
  @attr limit;
  @attr page;
  @attr sort;
  @attr total;
  @attr pages;
}

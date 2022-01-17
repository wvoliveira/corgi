import Model, { attr } from '@ember-data/model';

export default class LinkModel extends Model {
  @attr create_at;
  @attr updated_at;
  @attr domain;
  @attr keyword;
  @attr url;
  @attr title;
  @attr active;
}

import Model, { attr } from '@ember-data/model';

export default class Link extends Model {
  //@attr('string') id;
  @attr('string') created_at;
  @attr('string') updated_at;

  @attr('string') domain;
  @attr('string') keyword;
  @attr('string') url;
  @attr('string') title;
  @attr('string') active;
}

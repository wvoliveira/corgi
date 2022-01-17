import Controller from '@ember/controller';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';

export default class IndexController extends Controller {
  @service session;
  @service router;

  queryParams = ['page', 'limit'];
  page = 1;
  limit = 10;
}

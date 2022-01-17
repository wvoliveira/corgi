import Controller from '@ember/controller';
import { tracked } from '@glimmer/tracking';

export default class LinksController extends Controller {
  queryParams = ['page', 'limit'];

  @tracked page = 1;
  @tracked limit = 10;

  @tracked model;

  get filteredLinks() {
    let page = this.page;
    let limit = this.limit;

    let links = this.model;

    console.log('links');
    console.log(links);

    return links;
  }
}

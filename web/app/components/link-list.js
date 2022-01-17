import Component from '@glimmer/component';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';

export default class LinkListComponent extends Component {
  @service store;
  @service session;
  @service router;

  @tracked page = 1;
  @tracked limit = 10;

  @tracked pagination = {
    back: false,
    backPage: 1,
    next: false,
    nextPage: 2,
  };

  @tracked nextPage = false;

  @tracked links = {};
  @tracked isLoading = false;

  constructor() {
    super(...arguments);
    this.loadLinks();
  }

  async loadLinks() {
    console.log('this.args');
    console.log(this.args.query);

    console.log('page: ' + this.args.query.page);
    console.log('limit: ' + this.args.query.limit);

    if (this.args.query.page) {
      this.page = this.args.query.page;
    }
    if (this.args.query.limit) {
      this.limit = this.args.query.limit;
    }

    this.isLoading = true;

    let allLinks;
    await this.store
      .query('link', { page: this.page, limit: this.limit })
      .then(function (response) {
        allLinks = response;
      });
    this.links = allLinks;

    console.log(this.links.meta);
    this.isLoading = false;

    if (this.page > 1) {
      this.pagination.back = true;
    }
    if (this.links.meta.pages > 1) {
      this.pagination.next = true;
    }
  }

  @action
  translation() {
    this.router.transitionTo('/');
  }
}

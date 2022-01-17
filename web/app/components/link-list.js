import Component from '@glimmer/component';
import {tracked} from '@glimmer/tracking';
import {inject as service} from '@ember/service';

export default class LinkListComponent extends Component {
  @service store;
  @service session;

  @tracked page = 1;
  @tracked limit = 10;

  @tracked links = {};
  @tracked isLoading = false;

  constructor() {
    super(...arguments);
    this.loadLinks();
  }

  async loadLinks() {
    this.isLoading = true;
    let allLinks;

    await this.store
      .query('link', { offset: this.page, limit: this.limit })
      .then(function (response) {
        allLinks = response;

        console.log('response');
        console.log(response);
        console.log(response.get('content'));
      });
    this.isLoading = false;

    this.links = allLinks;
    console.log('this.links');
    console.log(this.links);
  }
}

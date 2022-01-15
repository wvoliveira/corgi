import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';

export default class LinkListComponent extends Component {
  @service session;

  @tracked links = {};
  @tracked isLoading = false;

  constructor() {
    super(...arguments);
    this.loadLinks();
  }

  async loadLinks() {
    this.isLoading = true;
    this.links = await this.session.fetch('/api/v1/links/');
    this.isLoading = false;

    console.log(this.links);
  }
}

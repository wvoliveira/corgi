import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class LinkListComponent extends Component {
  @service session;
  @service store;

  @tracked isRunning = null;
  @tracked links = [];

  constructor() {
    super(...arguments);
    this.loadLinks();
  }

  async loadLinks() {
    this.isRunning = true;

    this.links = await this.store.findAll('link').then(function (links) {
      console.log('links');
      console.log(links);
      return links;
    });
    console.log('this.links');
    console.log(this.links);
  }
}

import Component from '@glimmer/component';
import {tracked} from '@glimmer/tracking';
import {inject as service} from '@ember/service';

export default class LinkListComponent extends Component {
  @service store;
  @service session;

  @tracked links = {};
  @tracked isLoading = false;

  constructor() {
    super(...arguments);
    this.loadLinks();
  }

  loadLinks() {
    this.isLoading = true;
    let allLinks;

    this.store.findAll('link').then(function (links) {
      console.log('link');
      console.log(links);
      allLinks = links;
    });
    this.isLoading = false;

    this.links = allLinks;
  }
}

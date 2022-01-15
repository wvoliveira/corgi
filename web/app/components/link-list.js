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

  async loadLinks() {
    this.isLoading = true;
    this.links = await this.store.findAll('link').then(function (link) {
      console.log('link: ' + link);
      let data = link.get('data');
      console.log('data: ' + data);
      return data;
    });

    console.log('this.links ' + this.links);
    this.isLoading = false;
    console.log('this.links ' + this.links);
  }
}

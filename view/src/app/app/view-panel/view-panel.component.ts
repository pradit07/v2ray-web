import { Component, OnInit, Input, OnDestroy } from '@angular/core';
import { Panel, Element } from '../view/source';
import { HttpClient } from '@angular/common/http';
import { ServerAPI } from 'src/app/core/core/api';
import { ToasterService } from 'angular2-toaster';
import { I18nService } from 'src/app/core/i18n/i18n.service';
import { Utils } from 'src/app/core/utils';
import { isArray } from 'util';
@Component({
  selector: 'app-view-panel',
  templateUrl: './view-panel.component.html',
  styleUrls: ['./view-panel.component.scss']
})
export class ViewPanelComponent implements OnInit, OnDestroy {
  constructor(private httpClient: HttpClient,
    private toasterService: ToasterService,
    private i18nService: I18nService,
  ) { }
  @Input('panel')
  panel: Panel
  private _disabled = false
  get disabled(): boolean {
    return this._disabled
  }
  private _closed = false
  ngOnInit(): void {
  }
  ngOnDestroy() {
    this._closed = true
  }
  onClickTest() {
    console.log('test')
  }
  onClickAdd() {
    console.log('add')
  }
  onClickUpdate() {
    this._disabled = true
    this.httpClient.post<Array<any>>(ServerAPI.proxy.update, {
      id: this.panel.id,
    }).toPromise().then((data) => {
      if (this._closed) {
        return
      }
      this.panel.source.slice(0, this.panel.source.length)
      if (isArray(data) && data.length > 0) {
        for (let i = 0; i < data.length; i++) {
          this.panel.source.push(new Element(data[i]))
        }
        this.panel.source.sort(Element.compare)
      }
    }, (e) => {
      if (this._closed) {
        return
      }
      console.warn(e)
      this.toasterService.pop('error',
        this.i18nService.get('error'),
        Utils.resolveError(e),
      )
    }).finally(() => {
      this._disabled = false
    })
  }
}
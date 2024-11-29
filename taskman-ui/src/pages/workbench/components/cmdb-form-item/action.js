/*
 * @Author: wanghao7717 792974788@qq.com
 * @Date: 2024-11-19 17:02:46
 * @LastEditors: wanghao7717 792974788@qq.com
 * @LastEditTime: 2024-11-19 17:03:11
 * @FilePath: \wecube-plugins-taskman\taskman-ui\src\pages\workbench\components\cmdb-form-item\action.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
export const components = {
  number: {
    component: 'Input',
    type: 'number'
  },
  datetime: {
    component: 'DatePicker',
    type: 'datetimerange'
  },
  text: {
    component: 'Input',
    type: 'text'
  },
  select: {
    component: 'WeCMDBSelect',
    options: []
  },
  ref: {
    component: 'WeCMDBRefSelect',
    highlightRow: true
  },
  extRef: {
    component: 'WeCMDBRefSelect',
    highlightRow: true
  },
  multiSelect: {
    component: 'WeCMDBSelect',
    options: []
  },
  multiRef: {
    component: 'WeCMDBRefSelect'
  },
  textArea: {
    component: 'Input',
    type: 'text'
  },
  password: {
    component: 'WeCMDBCIPassword'
  },
  diffVariable: {
    component: 'WeCMDBDiffVariable'
  }
}

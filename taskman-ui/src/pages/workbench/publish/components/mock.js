export default {
  request: {
    id: '20240119-000028',
    name: '刘超测试-人工任务-240119180258',
    form: '65aa48d2975d1b1e',
    requestTemplate: '65aa48a6c2b72da4',
    requestTemplateName: '',
    procInstanceId: '640',
    procInstanceKey: 'u1NMDEpQ4Zjq',
    reporter: 'admin',
    handler: 'admin',
    reportTime: '2024-01-19 18:04:45',
    status: 'InProgress',
    customFormCache: '',
    result: '',
    expireTime: '2024-01-25 18:07:13',
    expectTime: '2024-01-25 18:04:45',
    confirmTime: '2024-01-19 18:07:13',
    createdBy: 'admin',
    createdTime: '2024-01-19 18:02:58',
    updatedBy: 'admin',
    updatedTime: '2024-01-19 18:08:17',
    delFlag: 0,
    handleRoles: null,
    parent: '',
    completedTime: '',
    rollbackDesc: ''
  },
  task: [
    {
      id: '65aa49d1d62fd615',
      name: '审批1',
      description: '',
      handler: 'admin',
      nextOption: '同意,拒绝',
      choseOption: '同意',
      createdBy: 'system',
      createdTime: '2024-01-19 18:07:13',
      updatedBy: 'admin',
      updatedTime: '2024-01-19 18:07:45',
      delFlag: '0',
      operationOptions: null,
      expireTime: '2024-01-21 18:07:13',
      notifyCount: 0,
      templateType: 0,
      type: 'approve',
      sort: 0,
      taskResult: '',
      confirmResult: '',
      editable: false,
      taskHandleList: [
        {
          id: '65aa49d1d62fd616',
          taskHandleTemplate: '65e092c7e4cdf9a1',
          task: '65aa49d1d62fd615',
          role: 'SUPER_ADMIN',
          handler: 'admin',
          handlerType: 'template',
          handleResult: '同意',
          resultDesc: '11111111111111',
          parentId: '',
          changeReason: '',
          createdTime: '2024-03-04 17:28:33',
          updatedTime: '2024-03-04 17:28:36',
          sort: 0
        },
        {
          id: '65aa49d1d62fd617',
          taskHandleTemplate: '65e092c7e4cdf9a1',
          task: '65aa49d1d62fd615',
          role: 'SUPER_ADMIN',
          handler: 'admin',
          handlerType: 'template',
          handleResult: '拒绝',
          resultDesc: '222222222222222',
          parentId: '',
          changeReason: '',
          createdTime: '2024-03-04 17:28:33',
          updatedTime: '2024-03-04 17:28:36',
          sort: 1
        }
      ],
      nextOptions: ['同意', '拒绝'],
      attachFiles: [
        {
          id: '65e6eb4b4994e018',
          name: '开发五室-转正答辩-社招-王浩(前端开发).pptx',
          s3BucketName: '',
          s3KeyName: '',
          delFlag: 0,
          request: '',
          task: ''
        }
      ],
      handleMode: 'all',
      formData: [
        {
          packageName: 'wecmdb',
          entity: 'app_instance',
          formTemplateId: '',
          itemGroup: 'wecmdb:app_instance',
          itemGroupName: 'wecmdb:app_instance',
          itemGroupType: 'workflow',
          itemGroupRule: 'new',
          title: [
            {
              id: '65d83c882910fd70',
              name: '新增测试',
              description: '',
              itemGroupId: '65d83c88a7bb9266',
              itemGroup: 'wecmdb:app_instance',
              itemGroupType: 'workflow',
              itemGroupName: 'wecmdb:app_instance',
              ItemGroupSort: 1,
              itemGroupRule: 'new',
              defaultValue: '1212',
              sort: 0,
              packageName: '',
              entity: '',
              attrDefId: '',
              attrDefName: '',
              attrDefDataType: '',
              elementType: 'calculate',
              title: '',
              width: 0,
              refPackageName: '',
              refEntity: '',
              dataOptions: '',
              required: '',
              regular: '',
              isEdit: '',
              isView: '',
              isOutput: '',
              inDisplayName: '',
              isRefInside: '',
              multiple: '',
              defaultClear: '',
              copyId: '',
              selectList: [],
              active: false
            },
            {
              id: '65d83c88510d9ee3',
              name: 'create_time',
              description: '创建时间',
              itemGroupId: '65d83c88a7bb9266',
              itemGroup: 'wecmdb:app_instance',
              itemGroupType: 'workflow',
              itemGroupName: 'wecmdb:app_instance',
              ItemGroupSort: 1,
              itemGroupRule: 'new',
              defaultValue: '',
              sort: 0,
              packageName: 'wecmdb',
              entity: 'app_instance',
              attrDefId: 'wecmdb:app_instance:create_time',
              attrDefName: 'create_time',
              attrDefDataType: 'str',
              elementType: 'input',
              title: '创建时间',
              width: 24,
              refPackageName: 'wecmdb',
              refEntity: 'create_time',
              dataOptions: '',
              required: 'no',
              regular: '',
              isEdit: 'yes',
              isView: 'yes',
              isOutput: 'no',
              inDisplayName: 'no',
              isRefInside: 'no',
              multiple: 'N',
              defaultClear: 'no',
              copyId: '',
              selectList: [],
              active: false
            },
            {
              id: '65d83c8869a3d4b3',
              name: 'asset_id',
              description: '资产ID',
              itemGroupId: '65d83c88a7bb9266',
              itemGroup: 'wecmdb:app_instance',
              itemGroupType: 'workflow',
              itemGroupName: 'wecmdb:app_instance',
              ItemGroupSort: 1,
              itemGroupRule: 'new',
              defaultValue: '',
              sort: 0,
              packageName: 'wecmdb',
              entity: 'app_instance',
              attrDefId: 'wecmdb:app_instance:asset_id',
              attrDefName: 'asset_id',
              attrDefDataType: 'str',
              elementType: 'input',
              title: '资产ID',
              width: 24,
              refPackageName: 'wecmdb',
              refEntity: 'asset_id',
              dataOptions: '',
              required: 'no',
              regular: '',
              isEdit: 'yes',
              isView: 'yes',
              isOutput: 'no',
              inDisplayName: 'no',
              isRefInside: 'no',
              multiple: 'N',
              defaultClear: 'no',
              copyId: '',
              selectList: [],
              active: false
            },
            {
              id: '65d83c88739d682c',
              name: '测试1111',
              description: '',
              itemGroupId: '65d83c88a7bb9266',
              itemGroup: 'wecmdb:app_instance',
              itemGroupType: 'workflow',
              itemGroupName: 'wecmdb:app_instance',
              ItemGroupSort: 1,
              itemGroupRule: 'new',
              defaultValue: '1212544545',
              sort: 0,
              packageName: '',
              entity: '',
              attrDefId: '',
              attrDefName: '',
              attrDefDataType: '',
              elementType: 'calculate',
              title: '',
              width: 0,
              refPackageName: '',
              refEntity: '',
              dataOptions: '',
              required: '',
              regular: '',
              isEdit: '',
              isView: '',
              isOutput: '',
              inDisplayName: '',
              isRefInside: '',
              multiple: '',
              defaultClear: '',
              copyId: '',
              selectList: [],
              active: false
            },
            {
              id: '65d83c88d2347f96',
              name: 'confirm_time',
              description: '确认时间',
              itemGroupId: '65d83c88a7bb9266',
              itemGroup: 'wecmdb:app_instance',
              itemGroupType: 'workflow',
              itemGroupName: 'wecmdb:app_instance',
              ItemGroupSort: 1,
              itemGroupRule: 'new',
              defaultValue: '',
              sort: 0,
              packageName: 'wecmdb',
              entity: 'app_instance',
              attrDefId: 'wecmdb:app_instance:confirm_time',
              attrDefName: 'confirm_time',
              attrDefDataType: 'str',
              elementType: 'input',
              title: '确认时间',
              width: 24,
              refPackageName: 'wecmdb',
              refEntity: 'confirm_time',
              dataOptions: '',
              required: 'no',
              regular: '',
              isEdit: 'yes',
              isView: 'yes',
              isOutput: 'no',
              inDisplayName: 'no',
              isRefInside: 'no',
              multiple: 'N',
              defaultClear: 'no',
              copyId: '',
              selectList: [],
              active: false
            }
          ],
          value: [
            {
              packageName: 'wecmdb',
              entityName: 'app_instance',
              dataId: '',
              displayName: '',
              fullDataId: '',
              id: 'tmp__65dc9321d126c84e',
              entityData: {
                _id: 'tmp__65dc9321d126c84e',
                asset_id: '的撒大苏打阿萨',
                confirm_time: '打撒大撒大撒大撒',
                create_time: '1122の2饿啊的撒旦撒旦撒旦',
                新增测试: '',
                测试1111: ''
              },
              previousIds: [],
              succeedingIds: [],
              entityDataOp: 'create'
            }
          ]
        }
      ]
    },
    {
      id: '65aa49f1d6a41b28',
      name: '任务1',
      description: '',
      form: '65aa49f174a48fdc',
      attachFile: '',
      status: 'implement_process',
      version: '',
      request: '20240119-000028',
      parent: '',
      taskTemplate: '65aa48a6f99fa65d',
      packageName: '',
      entityName: '',
      procDefId: 'u1NLiQcQ4YXk',
      procDefKey: 'wecube1705657726396',
      procDefName: '',
      nodeDefId: 'u1NLiQBQ4Z1T',
      nodeName: '任务1',
      callbackUrl: '/platform/v1/process/instances/callback',
      callbackParameter: 'SUTN-u1NMLZYQ4ZpZ',
      callbackData:
        '{"resultCode":"同意","resultMessage":"","results":{"requestId":"7e642e2131ec4ad6b73753af318d5c84","outputs":[{"callbackParameter":"SUTN-u1NMLZYQ4ZpZ","comment":"都是否","taskFormOutput":"{\\"formMetaId\\":\\"65aa48a6c8b36914\\",\\"procDefId\\":\\"u1NLiQcQ4YXk\\",\\"procDefKey\\":\\"wecube1705657726396\\",\\"procInstId\\":640,\\"procInstKey\\":\\"u1NMDEpQ4Zjq\\",\\"taskNodeDefId\\":\\"u1NLiQBQ4Z1T\\",\\"taskNodeInstId\\":4010,\\"formDataEntities\\":[]}","errorCode":"0","errorMessage":""}]}}',
      emergency: 0,
      result: '都是否',
      cache:
        '{"formMetaId":"65aa48a6c8b36914","procDefId":"u1NLiQcQ4YXk","procDefKey":"wecube1705657726396","procInstId":640,"procInstKey":"u1NMDEpQ4Zjq","taskNodeDefId":"u1NLiQBQ4Z1T","taskNodeInstId":4010,"formDataEntities":[]}',
      callbackRequestId: '7e642e2131ec4ad6b73753af318d5c84',
      reporter: '',
      reportTime: '2024-01-19 18:07:45',
      reportRole: '',
      handler: 'admin',
      nextOption: '同意,拒绝',
      choseOption: '同意',
      createdBy: 'system',
      createdTime: '2024-01-19 18:07:45',
      updatedBy: 'admin',
      updatedTime: '2024-01-19 18:08:16',
      delFlag: '0',
      operationOptions: null,
      expireTime: '2024-01-21 18:07:45',
      notifyCount: 0,
      templateType: 0,
      type: 'implement_process',
      sort: 0,
      taskResult: '',
      confirmResult: '',
      editable: true,
      taskHandleList: [
        {
          id: '65aa49d1d62fd616',
          taskHandleTemplate: '65e092c7e4cdf9a1',
          task: '65aa49d1d62fd615',
          role: 'SUPER_ADMIN',
          handler: 'admin',
          handlerType: 'template',
          handleResult: '同意',
          resultDesc: '1111111111111111111',
          parentId: '',
          changeReason: '',
          createdTime: '2024-03-04 17:28:33',
          updatedTime: '2024-03-04 17:28:36',
          sort: 0
        },
        {
          id: '65aa49d1d62fd617',
          taskHandleTemplate: '65e092c7e4cdf9a1',
          task: '65aa49d1d62fd615',
          role: 'SUPER_ADMIN',
          handler: 'admin',
          handlerType: 'template',
          handleResult: '',
          resultDesc: '',
          parentId: '',
          changeReason: '',
          createdTime: '2024-03-04 17:28:33',
          updatedTime: '2024-03-04 17:28:36',
          sort: 1
        }
      ],
      nextOptions: ['同意', '拒绝'],
      attachFiles: [
        {
          id: '65e6eb307258118d',
          name: '参数.txt',
          s3BucketName: '',
          s3KeyName: '',
          delFlag: 0,
          request: '',
          task: ''
        }
      ],
      handleMode: 'any',
      formData: [
        {
          packageName: 'wecmdb',
          entity: 'app_instance',
          formTemplateId: '',
          itemGroup: 'wecmdb:app_instance',
          itemGroupName: 'wecmdb:app_instance',
          itemGroupType: 'workflow',
          itemGroupRule: 'new',
          title: [
            {
              id: '65d83c882910fd70',
              name: '新增测试',
              description: '',
              itemGroupId: '65d83c88a7bb9266',
              itemGroup: 'wecmdb:app_instance',
              itemGroupType: 'workflow',
              itemGroupName: 'wecmdb:app_instance',
              ItemGroupSort: 1,
              itemGroupRule: 'new',
              defaultValue: '1212',
              sort: 0,
              packageName: '',
              entity: '',
              attrDefId: '',
              attrDefName: '',
              attrDefDataType: '',
              elementType: 'calculate',
              title: '',
              width: 0,
              refPackageName: '',
              refEntity: '',
              dataOptions: '',
              required: '',
              regular: '',
              isEdit: '',
              isView: '',
              isOutput: '',
              inDisplayName: '',
              isRefInside: '',
              multiple: '',
              defaultClear: '',
              copyId: '',
              selectList: [],
              active: false
            },
            {
              id: '65d83c88510d9ee3',
              name: 'create_time',
              description: '创建时间',
              itemGroupId: '65d83c88a7bb9266',
              itemGroup: 'wecmdb:app_instance',
              itemGroupType: 'workflow',
              itemGroupName: 'wecmdb:app_instance',
              ItemGroupSort: 1,
              itemGroupRule: 'new',
              defaultValue: '',
              sort: 0,
              packageName: 'wecmdb',
              entity: 'app_instance',
              attrDefId: 'wecmdb:app_instance:create_time',
              attrDefName: 'create_time',
              attrDefDataType: 'str',
              elementType: 'input',
              title: '创建时间',
              width: 24,
              refPackageName: 'wecmdb',
              refEntity: 'create_time',
              dataOptions: '',
              required: 'no',
              regular: '',
              isEdit: 'yes',
              isView: 'yes',
              isOutput: 'no',
              inDisplayName: 'no',
              isRefInside: 'no',
              multiple: 'N',
              defaultClear: 'no',
              copyId: '',
              selectList: [],
              active: false
            },
            {
              id: '65d83c8869a3d4b3',
              name: 'asset_id',
              description: '资产ID',
              itemGroupId: '65d83c88a7bb9266',
              itemGroup: 'wecmdb:app_instance',
              itemGroupType: 'workflow',
              itemGroupName: 'wecmdb:app_instance',
              ItemGroupSort: 1,
              itemGroupRule: 'new',
              defaultValue: '',
              sort: 0,
              packageName: 'wecmdb',
              entity: 'app_instance',
              attrDefId: 'wecmdb:app_instance:asset_id',
              attrDefName: 'asset_id',
              attrDefDataType: 'str',
              elementType: 'input',
              title: '资产ID',
              width: 24,
              refPackageName: 'wecmdb',
              refEntity: 'asset_id',
              dataOptions: '',
              required: 'no',
              regular: '',
              isEdit: 'yes',
              isView: 'yes',
              isOutput: 'no',
              inDisplayName: 'no',
              isRefInside: 'no',
              multiple: 'N',
              defaultClear: 'no',
              copyId: '',
              selectList: [],
              active: false
            },
            {
              id: '65d83c88739d682c',
              name: '测试1111',
              description: '',
              itemGroupId: '65d83c88a7bb9266',
              itemGroup: 'wecmdb:app_instance',
              itemGroupType: 'workflow',
              itemGroupName: 'wecmdb:app_instance',
              ItemGroupSort: 1,
              itemGroupRule: 'new',
              defaultValue: '1212544545',
              sort: 0,
              packageName: '',
              entity: '',
              attrDefId: '',
              attrDefName: '',
              attrDefDataType: '',
              elementType: 'calculate',
              title: '',
              width: 0,
              refPackageName: '',
              refEntity: '',
              dataOptions: '',
              required: '',
              regular: '',
              isEdit: '',
              isView: '',
              isOutput: '',
              inDisplayName: '',
              isRefInside: '',
              multiple: '',
              defaultClear: '',
              copyId: '',
              selectList: [],
              active: false
            },
            {
              id: '65d83c88d2347f96',
              name: 'confirm_time',
              description: '确认时间',
              itemGroupId: '65d83c88a7bb9266',
              itemGroup: 'wecmdb:app_instance',
              itemGroupType: 'workflow',
              itemGroupName: 'wecmdb:app_instance',
              ItemGroupSort: 1,
              itemGroupRule: 'new',
              defaultValue: '',
              sort: 0,
              packageName: 'wecmdb',
              entity: 'app_instance',
              attrDefId: 'wecmdb:app_instance:confirm_time',
              attrDefName: 'confirm_time',
              attrDefDataType: 'str',
              elementType: 'input',
              title: '确认时间',
              width: 24,
              refPackageName: 'wecmdb',
              refEntity: 'confirm_time',
              dataOptions: '',
              required: 'no',
              regular: '',
              isEdit: 'yes',
              isView: 'yes',
              isOutput: 'no',
              inDisplayName: 'no',
              isRefInside: 'no',
              multiple: 'N',
              defaultClear: 'no',
              copyId: '',
              selectList: [],
              active: false
            }
          ],
          value: [
            {
              packageName: 'wecmdb',
              entityName: 'app_instance',
              dataId: '',
              displayName: '',
              fullDataId: '',
              id: 'tmp__65dc9321d126c84e',
              entityData: {
                _id: 'tmp__65dc9321d126c84e',
                asset_id: '的撒大苏打阿萨',
                confirm_time: '打撒大撒大撒大撒',
                create_time: '1122の2饿啊的撒旦撒旦撒旦',
                新增测试: '',
                测试1111: ''
              },
              previousIds: [],
              succeedingIds: [],
              entityDataOp: 'create'
            }
          ]
        }
      ]
    }
  ]
}
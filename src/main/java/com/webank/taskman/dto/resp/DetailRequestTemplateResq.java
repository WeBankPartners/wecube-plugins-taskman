package com.webank.taskman.dto.resp;

public class DetailRequestTemplateResq {
    private String id;

    private String requestTempGroup;

    private String procDefKey;

    private String procDefId;

    private String procDefName;

    private String name;

    private String version;

    private String tags;

    private String status;

    private DetilReuestTemplateFormResq detilReuestTemplateFormResq;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getRequestTempGroup() {
        return requestTempGroup;
    }

    public void setRequestTempGroup(String requestTempGroup) {
        this.requestTempGroup = requestTempGroup;
    }

    public String getProcDefKey() {
        return procDefKey;
    }

    public void setProcDefKey(String procDefKey) {
        this.procDefKey = procDefKey;
    }

    public String getProcDefId() {
        return procDefId;
    }

    public void setProcDefId(String procDefId) {
        this.procDefId = procDefId;
    }

    public String getProcDefName() {
        return procDefName;
    }

    public void setProcDefName(String procDefName) {
        this.procDefName = procDefName;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public String getTags() {
        return tags;
    }

    public void setTags(String tags) {
        this.tags = tags;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public DetilReuestTemplateFormResq getDetilReuestTemplateFormResq() {
        return detilReuestTemplateFormResq;
    }

    public void setDetilReuestTemplateFormResq(DetilReuestTemplateFormResq detilReuestTemplateFormResq) {
        this.detilReuestTemplateFormResq = detilReuestTemplateFormResq;
    }
}

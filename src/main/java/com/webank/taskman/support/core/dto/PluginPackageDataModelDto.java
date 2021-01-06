package com.webank.taskman.support.core.dto;


import java.util.ArrayList;
import java.util.List;
import java.util.Set;
import java.util.StringJoiner;

public class PluginPackageDataModelDto {

    public static enum Source {
        PLUGIN_PACKAGE, DATA_MODEL_ENDPOINT
    }
    private String id;

    private Integer version;

    private String packageName;

    private boolean isDynamic;

    private String updatePath;

    private String updateMethod;

    private String updateSource;

    private Long updateTime;

    private List<PluginPackageEntityDto> entities = new ArrayList<>();

    public PluginPackageDataModelDto() {
    }

    public PluginPackageDataModelDto(String id, Integer version, String packageName, boolean isDynamic, String updatePath, String updateMethod, String updateSource, Long updateTime, List<PluginPackageEntityDto> entities) {
        this.id = id;
        this.version = version;
        this.packageName = packageName;
        this.isDynamic = isDynamic;
        this.updatePath = updatePath;
        this.updateMethod = updateMethod;
        this.updateSource = updateSource;
        this.updateTime = updateTime;
        this.entities = entities;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public Integer getVersion() {
        return version;
    }

    public void setVersion(Integer version) {
        this.version = version;
    }

    public String getPackageName() {
        return packageName;
    }

    public void setPackageName(String packageName) {
        this.packageName = packageName;
    }

    public boolean isDynamic() {
        return isDynamic;
    }

    public void setDynamic(boolean dynamic) {
        isDynamic = dynamic;
    }

    public String getUpdatePath() {
        return updatePath;
    }

    public void setUpdatePath(String updatePath) {
        this.updatePath = updatePath;
    }

    public String getUpdateMethod() {
        return updateMethod;
    }

    public void setUpdateMethod(String updateMethod) {
        this.updateMethod = updateMethod;
    }

    public String getUpdateSource() {
        return updateSource;
    }

    public void setUpdateSource(String updateSource) {
        this.updateSource = updateSource;
    }

    public Long getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(Long updateTime) {
        this.updateTime = updateTime;
    }

    public List<PluginPackageEntityDto> getEntities() {
        return entities;
    }

    public PluginPackageDataModelDto setEntities(List<PluginPackageEntityDto> entities) {
        this.entities = entities;
        return this;
    }

    @Override
    public String toString() {
        return new StringJoiner(", ", PluginPackageDataModelDto.class.getSimpleName() + "[", "]")
                .add("id='" + id + "'")
                .add("version=" + version)
                .add("packageName='" + packageName + "'")
                .add("isDynamic=" + isDynamic)
                .add("updatePath='" + updatePath + "'")
                .add("updateMethod='" + updateMethod + "'")
                .add("updateSource='" + updateSource + "'")
                .add("updateTime=" + updateTime)
                .add("entities=" + entities)
                .toString();
    }
}

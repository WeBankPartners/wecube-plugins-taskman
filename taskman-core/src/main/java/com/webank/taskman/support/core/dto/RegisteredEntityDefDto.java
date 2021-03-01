package com.webank.taskman.support.core.dto;

import java.util.ArrayList;
import java.util.List;

public class RegisteredEntityDefDto {
    private String id;
    private String packageName;
    private String name;
    private String displayName;
    private String description;

    private List<RegisteredEntityAttrDefDto> attributes = new ArrayList<>();

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getPackageName() {
        return packageName;
    }

    public void setPackageName(String packageName) {
        this.packageName = packageName;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDisplayName() {
        return displayName;
    }

    public void setDisplayName(String displayName) {
        this.displayName = displayName;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public List<RegisteredEntityAttrDefDto> getAttributes() {
        return attributes;
    }

    public void setAttributes(List<RegisteredEntityAttrDefDto> attributes) {
        this.attributes = attributes;
    }

    @Override
    public String toString() {
        StringBuilder builder = new StringBuilder();
        builder.append("RegisteredEntityDefDto [id=");
        builder.append(id);
        builder.append(", packageName=");
        builder.append(packageName);
        builder.append(", name=");
        builder.append(name);
        builder.append(", displayName=");
        builder.append(displayName);
        builder.append(", description=");
        builder.append(description);
        builder.append(", attributes=");
        builder.append(attributes);
        builder.append("]");
        return builder.toString();
    }

    public static class RegisteredEntityAttrDefDto {
        private String id;
        private String name;
        private String description;
        private String dataType;
        private boolean mandatory = false;

        private String refPackageName;
        private String refEntityName;
        private String refAttrName;

        public String getId() {
            return id;
        }

        public void setId(String id) {
            this.id = id;
        }

        public String getName() {
            return name;
        }

        public void setName(String name) {
            this.name = name;
        }

        public String getDescription() {
            return description;
        }

        public void setDescription(String description) {
            this.description = description;
        }

        public String getDataType() {
            return dataType;
        }

        public void setDataType(String dataType) {
            this.dataType = dataType;
        }

        public boolean isMandatory() {
            return mandatory;
        }

        public void setMandatory(boolean mandatory) {
            this.mandatory = mandatory;
        }

        public String getRefPackageName() {
            return refPackageName;
        }

        public void setRefPackageName(String refPackageName) {
            this.refPackageName = refPackageName;
        }

        public String getRefEntityName() {
            return refEntityName;
        }

        public void setRefEntityName(String refEntityName) {
            this.refEntityName = refEntityName;
        }

        public String getRefAttrName() {
            return refAttrName;
        }

        public void setRefAttrName(String refAttrName) {
            this.refAttrName = refAttrName;
        }

        @Override
        public String toString() {
            StringBuilder builder = new StringBuilder();
            builder.append("RegisteredEntityAttrDefDto [id=");
            builder.append(id);
            builder.append(", name=");
            builder.append(name);
            builder.append(", description=");
            builder.append(description);
            builder.append(", dataType=");
            builder.append(dataType);
            builder.append(", mandatory=");
            builder.append(mandatory);
            builder.append(", refPackageName=");
            builder.append(refPackageName);
            builder.append(", refEntityName=");
            builder.append(refEntityName);
            builder.append(", refAttrName=");
            builder.append(refAttrName);
            builder.append("]");
            return builder.toString();
        }
    }
}

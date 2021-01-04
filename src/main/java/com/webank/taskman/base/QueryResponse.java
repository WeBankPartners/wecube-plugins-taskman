package com.webank.taskman.base;

import com.baomidou.mybatisplus.core.metadata.IPage;
import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonInclude.Include;

import java.util.LinkedList;
import java.util.List;

@JsonInclude(Include.NON_NULL)
public class QueryResponse<T> {
    private PageInfo pageInfo = new PageInfo();
    private List<T> contents = new LinkedList<>();

    public QueryResponse() {
    }

    public QueryResponse(Long totalRows, Long startIndex, Long pageSize, List<T> contents) {
        this.pageInfo = new PageInfo(totalRows,startIndex,pageSize);
        this.contents = contents;
    }
    public QueryResponse(PageInfo pageInfo, List<T> contents) {
        this.pageInfo = pageInfo;
        this.contents = contents;
    }

    public QueryResponse(IPage<T> iPage, List<T> contents) {
        this.pageInfo = new PageInfo(iPage.getTotal(),iPage.getCurrent(),iPage.getSize());
        this.contents = contents;
    }


    public PageInfo getPageInfo() {
        return pageInfo;
    }

    public void setPageInfo(PageInfo pageInfo) {
        this.pageInfo = pageInfo;
    }

    public List<T> getContents() {
        return contents;
    }

    public void setContents(List<T> domainObjs) {
        this.contents = domainObjs;
    }

    public void addContent(T ciObj) {
        this.contents.add(ciObj);
    }

    public static class Pageable {
        private Long startIndex;
        private Long pageSize = 10000L;

        public Pageable() {
        }

        public Pageable(Long startIndex, Long pageSize) {
            this.startIndex = startIndex;
            this.pageSize = pageSize;
        }

        public Long getPageSize() {
            return pageSize;
        }
        public void setPageSize(Long pageSize) {
            this.pageSize = pageSize;
        }
        public Long getStartIndex() {
            return startIndex;
        }
        public void setStartIndex(Long startIndex) {
            this.startIndex = startIndex;
        }
    }

    public static class PageInfo extends Pageable {
        private Long totalRows;

        public PageInfo() {
            this.totalRows = 0L;
        }

        public PageInfo(Long totalRows, Long startIndex, Long pageSize) {
            super(startIndex, pageSize);
            this.totalRows = totalRows;
        }

        public Long getTotalRows() {
            return totalRows;
        }

        public void setTotalRows(Long totalRows) {
            this.totalRows = totalRows;
        }
    }
}

package com.webank.taskman.utils;

import java.text.DateFormat;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;
import java.util.TimeZone;

import com.google.common.base.Strings;

public class DateUtils {

    public static Date convertToTimestamp(String value) throws Exception {
        if (Strings.isNullOrEmpty(value))
            return null;

        String timestampPattern = "yyyy-MM-dd HH:mm:ss";
        String datePattern = "yyyy-MM-dd";

        String parsePattern = null;
        Date date = null;
        if (value.length() == timestampPattern.length()) {
            parsePattern = timestampPattern;
        } else if (value.length() == datePattern.length()) {
            parsePattern = datePattern;
        }

        if (!Strings.isNullOrEmpty(parsePattern)) {
            SimpleDateFormat dateFmt = new SimpleDateFormat(parsePattern);
            dateFmt.setTimeZone(TimeZone.getTimeZone("UTC"));
            try {
                date = dateFmt.parse(value);
            } catch (ParseException e) {
                throw new Exception(String.format("Failed to parse date string [%s].", value), e);
            }
        } else {
            throw new Exception("Only support 'yyyy-MM-dd HH:mm:ss' and 'yyyy-MM-dd' for datetime.");
        }

        return date;
    }

    public static String formatDateToString(Date date) {
        if (date == null) {
            return "";
        }
        DateFormat df = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        df.setTimeZone(TimeZone.getDefault());
        return df.format(date);
    }

    public static Date formatStringToDate(String dateString) throws ParseException {
        DateFormat df = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
        return df.parse(dateString);
    }

    public static Date addDateMinute(Date date,String minute)  {
        Calendar calendar = Calendar.getInstance();
        calendar.setTime(date);
        calendar.add(Calendar.MINUTE, Integer.parseInt(minute));
        return calendar.getTime();
    }

}

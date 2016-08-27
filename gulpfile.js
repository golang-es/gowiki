'use strict';

var gulp        = require("gulp");
var sass        = require("gulp-sass");
var sourcemaps  = require('gulp-sourcemaps');

gulp.task('sass', function () {

    return gulp.src('./resources/sass/**/*.scss')
        .pipe(sourcemaps.init())
        .pipe(sass().on('error', sass.logError))
        .pipe(sourcemaps.write('./'))
        .pipe(gulp.dest('./public/css'));
});

gulp.task('watch', function () {
    gulp.watch('./resources/sass/**/*.scss', ['sass']);
})

gulp.task('default', ['watch']);

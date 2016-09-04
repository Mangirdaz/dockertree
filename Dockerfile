FROM centos
ADD dockertree /
ENV PORT 8080 
EXPOSE 8080
ADD server/dummydata.json /
CMD ["./dockertree", "-runmode", "server"]